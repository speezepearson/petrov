import React from 'react';

type TimerProps = {
    zeroTime: Date;
}
type TimerState = {
    secondsRemaining: number;
}

function twoDigits(x: number): string {
    return (x < 10) ? `0${x}` : `${x}`;
}

function formatSeconds(seconds: number): string {
    let dt, m, s, hund: number;
    const negative = seconds < 0;

    dt = Math.abs(seconds);
    [m, dt] = [Math.floor(dt / 60), dt % 60];
    [s, dt] = [Math.floor(dt), dt % 1];
    [hund, dt] = [Math.floor(100*dt), dt % .01];

    return `${negative ? '-' : ''}${twoDigits(m)}:${twoDigits(s)}${(!negative && seconds<15) ? "."+twoDigits(hund) : ""}`;
}

export class Timer extends React.Component<TimerProps, TimerState> {
    constructor(props: TimerProps) {
        super(props);
        this.state = {secondsRemaining: 0};
        this.tick();
    }
    private updaterId?: number;
    componentWillMount() {
        this.updaterId = setInterval(this.tick.bind(this), 10);
    }
    componentWillUnmount() {
        if (this.updaterId) {
            clearInterval(this.updaterId);
        }
    }

    render() {
        return formatSeconds(this.state.secondsRemaining);
    }

    tick() {
        let secondsRemaining = Math.min(this.state.secondsRemaining);
        this.setState({
            secondsRemaining: (this.props.zeroTime.getTime() - new Date().getTime()) / 1000,
        });
    }
}
