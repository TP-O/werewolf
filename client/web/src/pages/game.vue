<script setup lang="ts">
import { Game, Scene } from 'phaser'

defineOptions({
  name: 'GamePage',
})

class Example extends Scene {
  preload() {
    this.load.html('nameform', 'assets/text/loginform.html')
    this.load.image('pic', 'assets/img/turkey-1985086.jpg')
  }

  create() {
    this.add.image(400, 300, 'pic')

    const text = this.add.text(10, 10, 'Please login to play', { color: 'white', fontFamily: 'Arial', fontSize: '32px ' })

    const element = this.add.dom(400, 600).createFromCache('nameform')

    element.setPerspective(800)

    element.addListener('click')

    element.on('click', (event: any) => {
      if (event.target.name === 'loginButton') {
        const inputUsername = element.getChildByName('username')
        const inputPassword = element.getChildByName('password')

        console.log(inputUsername)

        //  Have they entered anything?
        if (1 === 1) {
          //  Turn off the click events
          this.events.removeListener('click')

          //  Tween the login form out
          this.tweens.add({ targets: element.rotate3d, x: 1, w: 90, duration: 3000, ease: 'Power3' })

          this.tweens.add({
            targets: element,
            scaleX: 2,
            scaleY: 2,
            y: 700,
            duration: 3000,
            ease: 'Power3',
            onComplete() {
              element.setVisible(false)
            },
          })

          //  Populate the text with whatever they typed in as the username!
          //   text.setText(`Welcome ${inputUsername.value}`)
          text.setText('hello em')
        }
        else {
          //  Flash the prompt
          this.tweens.add({ targets: text, alpha: 0.1, duration: 200, ease: 'Power3', yoyo: true })
        }
      }
    })

    this.tweens.add({
      targets: element,
      y: 300,
      duration: 3000,
      ease: 'Power3',
    })
  }
}

const config = {
  type: Phaser.AUTO,
  width: 800,
  height: 600,
  parent: 'phaser-example',
  dom: {
    createContainer: true,
  },
  scene: Example,
}

let game: Game | null = null

onMounted(() => {
  game = new Game(config)
})

onUnmounted(() => {
  game?.destroy(false)
})
</script>

<template>
  <div id="phaser-example" />
</template>
