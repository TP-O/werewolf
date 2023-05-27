<script setup lang="ts">
import log from 'loglevel'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import { CommEmitEvent } from '~/enums'
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

const isRoleSelectionShowed = ref(false)
function showRoleSelection() {
  isRoleSelectionShowed.value = true
}

const isRoleDetailsShowed = ref(false)
function showRoleDetails() {
  isRoleDetailsShowed.value = true
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

      <p text-xl>
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
        <div mb-2>Settings</div>

        <div grid="~ cols-2" mb-4 gap-4>
          <q-input
            outlined
            model-value="dsds"
            :disable="!isOwner"
            label="Password"
          />
          <q-input
            outlined
            :model-value="1"
            type="number"
            min="5"
            max="20"
            :disable="!isOwner"
            label="Capacity"
          />
          <q-input
            outlined
            :model-value="1"
            type="number"
            min="20"
            max="60"
            :disable="!isOwner"
            label="Turn duration"
          />
          <q-input
            outlined
            :model-value="1"
            type="number"
            min="40"
            max="360"
            :disable="!isOwner"
            label="Discussion duration"
          />
        </div>

        <q-btn
          color="secondary"
          :label="isOwner ? 'Add role' : 'View roles'"
          mb-4
          w-full
          outline
          @click="showRoleSelection"
        />

        <div grid="~ cols-1" gap-4 overflow-y-scroll>
          <q-card v-for="i in 10" :key="i" class="my-card" flat bordered>
            <q-card-section horizontal justify-between>
              <q-item>
                <q-item-section avatar>
                  <q-avatar>
                    <img src="https://cdn.quasar.dev/img/avatar2.jpg" />
                  </q-avatar>
                </q-item-section>

                <q-item-section>
                  <q-item-label>Title</q-item-label>
                  <q-item-label caption>Subhead</q-item-label>
                </q-item-section>
              </q-item>

              <q-card-actions>
                <q-btn flat round color="red" icon="favorite" />
                <q-btn flat round color="red" icon="favorite" />
              </q-card-actions>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <div w="1/3" flex="~ col" border rounded p-1>
        <div mb-2>Map</div>

        <q-card>
          <q-img src="https://cdn.quasar.dev/img/parallax2.jpg">
            <div class="absolute-bottom">
              <div class="text-h6">Our Changing Planet</div>
              <div class="text-subtitle2">by John Doe</div>
            </div>
          </q-img>
        </q-card>

        <div mt-8 grid="~ cols-2" gap-4 overflow-y-scroll>
          <q-avatar v-for="i in 10" :key="i" square h-auto w-full>
            <img src="https://cdn.quasar.dev/img/parallax2.jpg" />
          </q-avatar>
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
        <q-card
          v-for="i in 10"
          :key="i"
          class="my-card"
          flat
          bordered
          cursor-pointer
          @click="showRoleDetails"
        >
          <q-card-section horizontal justify-between>
            <q-item>
              <q-item-section avatar>
                <q-avatar>
                  <img src="https://cdn.quasar.dev/img/avatar2.jpg" />
                </q-avatar>
              </q-item-section>

              <q-item-section>
                <q-item-label>Title</q-item-label>
                <q-item-label caption>Subhead</q-item-label>
              </q-item-section>
            </q-item>

            <q-card-actions>
              <q-btn flat round color="red" icon="favorite" />
            </q-card-actions>
          </q-card-section>
        </q-card>
      </q-card-section>
    </q-card>
  </q-dialog>

  <q-dialog
    v-model="isRoleDetailsShowed"
    transition-show="rotate"
    transition-hide="rotate"
  >
    <q-card w="700px" max-w="80vw" h="500px" max-h="80vh">
      <q-img src="https://cdn.quasar.dev/img/parallax2.jpg" fit="contain" />

      <q-card-section>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod
        tempor incididunt ut labore et dolore magna aliqua.
      </q-card-section>
    </q-card>
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
