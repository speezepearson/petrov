import React from 'react';

type TimerProps = {
    currentTime: Date;
    zeroTime: Date;
    showHours?: boolean;
    showFractionBelow?: number;
}

function twoDigits(x: number): string {
    return (x < 10) ? `0${x}` : `${x}`;
}

export function formatSeconds(seconds: number, options: {showHours?: boolean, showFractionBelow?: number} = {}): string {
    const negative = seconds < 0;

    const absoluteSec = Math.abs(seconds);
    const [remainingSec, fractional] = [Math.floor(absoluteSec), absoluteSec - Math.floor(absoluteSec)];
    const [remainingMin, sec] = [Math.floor(remainingSec / 60), remainingSec % 60];
    const [hr, min] = [Math.floor(remainingMin / 60), remainingMin % 60];

    const segments: string[] = [
        negative ? '-' : '',
        options.showHours ? `${twoDigits(hr)}:${twoDigits(min)}:` : `${twoDigits(60*hr + min)}:`,
        twoDigits(sec),
        (options.showFractionBelow && (absoluteSec < options.showFractionBelow)) ? `.${twoDigits(Math.floor(100*fractional))}` : '',
    ]
    return segments.join('');
}

export function Timer(props: TimerProps) {
    return <span>
        {formatSeconds(
            (props.zeroTime.getTime() - props.currentTime.getTime()) / 1000,
            props,
        )}
    </span>;
}
