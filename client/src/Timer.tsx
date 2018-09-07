import React from 'react';

type TimerProps = {
    currentTime: Date;
    zeroTime: Date;
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

export function Timer(props: TimerProps) {
    return <span>
        {formatSeconds((props.zeroTime.getTime() - props.currentTime.getTime()) / 1000)}
    </span>;
}
