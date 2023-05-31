<script setup lang="ts">
import log from 'loglevel'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import { roles } from '~/constants'
import { CommEmitEvent, RoleId } from '~/enums'
import type { ResponseError } from '~/types'

defineOptions({
  name: 'WaitingRoomPage',
})

const route = useRoute()
const router = useRouter()
const $q = useQuasar()
const { player } = storeToRefs(usePlayerStore())
const { room, messages } = storeToRefs(useWaitingRoomStore())
const { join, leave, kick } = useWaitingRoomStore()

const roomId = String(route.params.id)
const isOwner = computed(() => room.value?.ownerId === player.value?.id)

async function leaveRoom() {
  $q.loading.show({
    message: 'Leaving...',
  })
  await leave(roomId)
}

function joinRoom(password?: string) {
  return join(roomId, password).catch(
    ({ data: { statusCode } }: ResponseError) => {
      if (statusCode === 400) {
        $q.dialog({
          title: 'Enter room password',
          prompt: {
            model: '',
            type: 'text',
          },
          cancel: true,
          persistent: true,
        }).onOk((password: string) => joinRoom(password))
      } else {
        router.push('/')

        if (statusCode === 404) {
          throw new Error('Room does not exist')
        } else {
          throw new Error('Unable to join this room')
        }
      }
    }
  )
}

const messageInput = ref('')
function sendMessage() {
  const data = {
    roomId,
    content: messageInput.value,
  }

  messageInput.value = ''
  log.info(`Send event ${CommEmitEvent.RoomMessage}:`, data)
  useCommSocket((socket) => socket.emit(CommEmitEvent.RoomMessage, data))
}

interface RoomSettings {
  password?: string
  capacity?: number
  gameSettings: {
    roleIds: RoleId[]
    requiredRoleIds: RoleId[]
    turnDuration: number
    discussionDuration: number
    mapId: number
  }
}

const roomSettings = reactive<RoomSettings>({
  capacity: 5,
  gameSettings: {
    roleIds: [RoleId.Villager, RoleId.Werewolf],
    requiredRoleIds: [RoleId.Villager, RoleId.Werewolf],
    turnDuration: 30,
    discussionDuration: 60,
    mapId: 1,
  },
})

const isRoleSelectionShowed = ref(false)
function showRoleSelection() {
  isRoleSelectionShowed.value = true
}

const isRoleDetailsShowed = ref(false)
const exploredRoleId = ref<RoleId>(1)
function showRoleDetails(id: number) {
  exploredRoleId.value = id
  isRoleDetailsShowed.value = true
}

function pickRole(id: number) {
  const i = roomSettings.gameSettings.roleIds.indexOf(id)
  if (i === -1) {
    roomSettings.gameSettings.roleIds.unshift(id)
  } else {
    roomSettings.gameSettings.roleIds.splice(i, 1)

    const j = roomSettings.gameSettings.requiredRoleIds.indexOf(id)
    if (j !== -1) {
      roomSettings.gameSettings.requiredRoleIds.splice(j, 1)
    }
  }
}

function pickRequiredRole(id: number) {
  const i = roomSettings.gameSettings.requiredRoleIds.indexOf(id)
  if (i === -1) {
    roomSettings.gameSettings.requiredRoleIds.unshift(id)

    if (!roomSettings.gameSettings.roleIds.includes(id)) {
      roomSettings.gameSettings.roleIds.unshift(id)
    }
  } else {
    roomSettings.gameSettings.requiredRoleIds.splice(i, 1)
  }
}

onBeforeMount(async () => {
  $q.loading.show({
    message: 'Loading room...',
  })

  if (!room.value || room.value?.id !== roomId) {
    await joinRoom()
  }

  $q.loading.hide()
})

const boxChat = ref<HTMLDivElement>()
onMounted(() => {
  if (boxChat.value) {
    const observer = new MutationObserver(() => {
      boxChat.value!.scrollTop = boxChat.value!.scrollHeight
    })
    observer.observe(boxChat.value, { childList: true })
  }
})

onUnmounted(() => {
  if ($q.loading.isActive) {
    $q.loading.hide()
  }
})
</script>

<template>
  <div h-full>
    <div h="1/12" grid="~ cols-[0.5fr_11fr_0.5fr]" py-4>
      <q-btn color="negative" label="Exit" @click="leaveRoom" />

      <p text="xl center">
        Room: <b>{{ roomId }}</b>
      </p>
    </div>

    <div h="11/12" flex gap-1>
      <div flex="~ col" gap-1 w="1/3">
        <div h="1/2" overflow-x-hidden overflow-y-scroll border rounded p-2>
          <div mb-2>Joined players</div>

          <div
            v-for="(memberId, i) of room?.memberIds"
            :key="i"
            flex="~  justify-between"
            my-2
            gap-2
            border
            p-2
          >
            <p
              flex="~ items-center"
              overflow-hidden
              text-ellipsis
              whitespace-nowrap
            >
              {{ memberId }}
            </p>
            <q-btn>
              <div i="carbon-overflow-menu-horizontal" />
              <q-menu min-w="100px" dense>
                <q-list>
                  <q-item
                    v-if="room?.ownerId === player?.username"
                    v-close-popup
                    clickable
                    @click="kick(memberId)"
                  >
                    <q-item-section>Kick</q-item-section>
                  </q-item>
                  <q-item v-close-popup clickable>
                    <q-item-section>Profile</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </div>
        </div>

        <div h="1/2" relative border rounded p-2>
          <div
            ref="boxChat"
            h="10/12"
            overflow-x-hidden
            overflow-y-scroll
            px-1
            text-left
          >
            <div v-for="(message, i) in messages" :key="i">
              <p mb-2 break-words>
                <b
                  v-if="message.senderId"
                  :class="
                    message.senderId === player?.username ? 'text-blue' : ''
                  "
                >
                  {{ `${message.senderId}:` }}
                </b>
                {{ message.content }}
              </p>
            </div>
          </div>

          <div absolute bottom-0 left-0 w-full bg-white p-2>
            <q-input
              v-model="messageInput"
              dense
              autogrow
              outlined
              placeholder="Say something..."
              h="max-[200px]"
              @keydown.enter.prevent="sendMessage"
            />
          </div>
        </div>
      </div>

      <div border rounded p-2 w="1/3" flex="~ col">
        <div mb-2>Picked roles</div>

        <q-btn
          color="secondary"
          :label="isOwner ? 'Add role' : 'View roles'"
          mb-4
          w-full
          outline
          @click="showRoleSelection"
        />

        <div grid="~ cols-1" gap-4 overflow-y-scroll>
          <RoleSummaryCard
            v-for="(role, i) of roomSettings.gameSettings.roleIds.map(
              (id) => roles[id]
            )"
            :key="i"
            :role="role"
            :picked="roomSettings.gameSettings.roleIds.includes(role.id)"
            :required="
              roomSettings.gameSettings.requiredRoleIds.includes(role.id)
            "
            :pick="pickRole"
            :mark-as-required="pickRequiredRole"
            @click="showRoleDetails(role.id)"
          />
        </div>
      </div>

      <div w="1/3" flex="~ col" relative border rounded p-1>
        <div mb-2>Settings</div>

        <div grid="~ cols-2" mb-4 gap-4>
          <q-input
            v-model="roomSettings.password"
            outlined
            :disable="!isOwner"
            label="Password"
          />
          <q-input
            v-model="roomSettings.capacity"
            outlined
            type="number"
            min="5"
            max="20"
            :disable="!isOwner"
            label="Capacity"
          />
          <q-input
            v-model="roomSettings.gameSettings.turnDuration"
            outlined
            type="number"
            min="20"
            max="60"
            :disable="!isOwner"
            label="Turn duration"
          />
          <q-input
            v-model="roomSettings.gameSettings.discussionDuration"
            outlined
            type="number"
            min="40"
            max="360"
            :disable="!isOwner"
            label="Discussion duration"
          />
        </div>

        <q-carousel
          v-model="roomSettings.gameSettings.mapId"
          animated
          arrows
          navigation
          infinite
        >
          <q-carousel-slide
            :name="1"
            img-src="https://cdn.quasar.dev/img/mountains.jpg"
          />
          <q-carousel-slide
            :name="2"
            img-src="https://cdn.quasar.dev/img/parallax1.jpg"
          />
          <q-carousel-slide
            :name="3"
            img-src="https://cdn.quasar.dev/img/parallax2.jpg"
          />
          <q-carousel-slide
            :name="4"
            img-src="https://cdn.quasar.dev/img/quasar.jpg"
          />
        </q-carousel>

        <div absolute bottom-0 left-0 w-full p-2>
          <q-btn color="secondary" label="Save" w-full outline />
        </div>
      </div>
    </div>
  </div>

  <q-dialog
    v-model="isRoleSelectionShowed"
    transition-show="rotate"
    transition-hide="rotate"
    full-height
  >
    <q-card w="700px" max-w="80vw">
      <q-card-section>
        <div class="text-h6">Roles</div>
      </q-card-section>

      <q-card-section flex="~ col" gap-4 pt-0>
        <RoleSummaryCard
          v-for="(role, i) of roles"
          :key="i"
          :role="role"
          :picked="roomSettings.gameSettings.roleIds.includes(role.id)"
          :required="
            roomSettings.gameSettings.requiredRoleIds.includes(role.id)
          "
          :pick="pickRole"
          :mark-as-required="pickRequiredRole"
          @click="showRoleDetails(role.id)"
        />
      </q-card-section>
    </q-card>
  </q-dialog>

  <q-dialog
    v-model="isRoleDetailsShowed"
    transition-show="rotate"
    transition-hide="rotate"
  >
    <RoleCard :role="roles[exploredRoleId]" />
  </q-dialog>
</template>

<style>
textarea {
  max-height: 100px;
}
</style>

<route lang="yaml">
meta:
  layout: client
</route>
