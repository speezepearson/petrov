import klaxonSound from './klaxon.mp3';
import missileSound from './launch.mp3'
import kaboomSound from './kaboom.mp3'

class Sound {

    private audio: HTMLAudioElement;
    private durationUpperBoundMs: number;
    private playing: boolean;

    constructor(soundFile: string, durationUpperBoundMs: number) {
        this.audio = new Audio(soundFile);
        this.durationUpperBoundMs = durationUpperBoundMs;
        this.playing = false;
    }

    play() {
        if (this.playing) return;
        this.playing = true;
        this.audio.play();
        window.setTimeout(() => this.reset(), this.durationUpperBoundMs);
    }

    reset() {
        this.audio.pause();
        this.audio.currentTime = 0;
        this.playing = false;
    }
}

export let KLAXON = new Sound(klaxonSound, 6000);
export let FWOOSH = new Sound(missileSound, 10000);
export let KABOOM = new Sound(kaboomSound, 30000);
