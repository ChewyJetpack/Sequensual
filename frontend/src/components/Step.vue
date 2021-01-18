<template>
    <div class="step">
        <div :class="[ 'step__indicator', { 'step__indicator--set': step.Trig.Active, 'step__indicator--active': step.Active }]" @click="toggleStep" />
        {{ step.Number + 1 }}
        <div class="step__controls">
            <div>Pitch:</div>
            <el-knob
                :value="pitch.value"
                :options="pitchKnobOpts"
                @input="changePitch($event.value)"
                size="xs"
            />
            <div>Octave:</div>
            <el-knob
                :value="pitch.octave"
                :options="[0, 1, 2, 3, 4, 5]"
                @input="changeOct($event)"
                size="xs"
            />

        </div>
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
    data() {
        return {
            notes: [
                "C",
                "C#",
                "D",
                "D#",
                "E",
                "F",
                "F#",
                "G",
                "G#",
                "A",
                "Bb",
                "B"
            ],
            pitch: {
                // '24' is the int value for C1 in MIDI
                value: 24,
                note: "C1",
                octave: 1
            }
        }
    },
    computed: {
        pitchKnobOpts() {
            let noteArr = [];
            for (let i = 0; i < this.notes.length; i++) {
                let opt = {
                    value: i + 24,
                    label: this.notes[i]
                }
                noteArr.push(opt)
            }
            return noteArr;
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
            this.step.Trig.Active = !this.step.Trig.Active;
            this.changeStep()
        },
        changeStep() {
            const number = this.step.Number;
            let status = this.step.Trig.Active;
            let note = this.pitch.value * this.pitch.octave;
            console.log(note)
            let velocity = 100;
            let length = 1;

            window.wails.Events.Emit("changeStep", {number, status, note, length, velocity});
        },
        changePitch(val) {
            this.pitch.value = val;
            this.changeStep()
        },
        changeOct(val) {
            this.pitch.octave = val;
            this.changeStep()
        }
    }
}
</script>

<style lang="scss" scoped>
.step {
    display: flex;
    flex-direction: column;
    align-items: center;
    
    &__indicator {
        height: 48px;
        width: 48px;
        background: crimson;
        border-radius: 50%;
        margin-bottom: 10px;
        cursor: pointer;

        &--active {
            background: orange;
        }

        &--set {
            background: yellowgreen;
        }
    }
}
</style>