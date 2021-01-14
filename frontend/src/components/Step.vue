<template>
    <div class="step">
        <div :class="[ 'step-indicator', { 'step-indicator--set': step.Trig.Active, 'step-indicator--active': step.Active }]" @click="toggleStep" />
        {{ step.Number + 1 }}
    </div>
</template>

<script>
export default {
    props: {
        step: {
            type: Object
        },
        callback: {
            type: Function
        }
    },
    mounted() {
        window.wails.Events.On("activeStep", stepNumber => {
            if (this.step.Number == stepNumber) {
                this.step.Active = true;
            } else {
                this.step.Active = false;
            }
        })
    },
    methods: {
        toggleStep() {
            const stepNumber = this.step.Number;
            let trigStatus = this.step.Trig.Active;
            trigStatus = !trigStatus;
            window.wails.Events.Emit("changeStep", {stepNumber, trigStatus});
            console.log('step toggled:', { number: stepNumber, status: trigStatus })
        }
    }
}
</script>

<style scoped>
.step {
    display: flex;
    flex-direction: column;
    align-items: center;
}
.step-indicator {
    height: 20px;
    width: 20px;
    background: crimson;
    border-radius: 50%;
    margin-bottom: 10px;
    cursor: pointer;
}

.step-indicator--active {
    background: orange;
}

.step-indicator--set {
    background: yellowgreen;
}
</style>