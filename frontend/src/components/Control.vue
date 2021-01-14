<template>
  <div class="container">
    <div class="transport">
      <a @click="play" class="transport__play">Play</a>
      <a @click="stop" class="transport__stop">Stop</a>
    </div>
    <div class="steps">
      <Step
        v-for="step in steps"
        :key="step.Number"
        :step="step"
        />
    </div>
  </div>
</template>

<script>
import Step from '@/components/Step'

export default {
  components: {
    Step
  },
  data() {
    return {
      steps: this.getSteps()
    }
  },
  mounted() {
    console.log(window)
    window.wails.Events.On("updateSteps", steps => {
      console.log(steps)
      if (steps) {
        this.steps = steps;
      }
    })
  },
  methods: {
    getSteps() {
      window.backend.Sequencer.GetSteps().then(result => {
        this.steps = result
      });
    },
    play() {
      window.wails.Events.Emit('play')
    },
    stop() {
      window.wails.Events.Emit('stop')
    }
    
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1 {
  margin-top: 2em;
  position: relative;
  min-height: 5rem;
  width: 100%;
}
a:hover {
  font-size: 1.7em;
  border-color: blue;
  background-color: blue;
  color: white;
  border: 3px solid white;
  border-radius: 10px;
  padding: 9px;
  cursor: pointer;
  transition: 500ms;
}
a {
  font-size: 1.7em;
  border-color: white;
  background-color: #121212;
  color: white;
  border: 3px solid white;
  border-radius: 10px;
  padding: 9px;
  cursor: pointer;
}

.transport {
  display: flex;
  justify-content: center;
}

.transport__play {
  margin-right: 10px;
}

.steps {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(16, 1fr);
  grid-gap: 10px;
  padding-top:40px;
}
</style>
