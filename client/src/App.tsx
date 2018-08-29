import jQuery from 'jquery';
import React from 'react';

import { Timer } from './Timer'
import {LaunchOrConcealButton} from "./LaunchOrConcealButton";

enum Phase {
    ENDED = "Ended",
    PRESTART = "PreStart",
    RUNNING = "Running",
    OVERTIME = "Overtime",
}

function nowPlus(seconds: number): Date {
    const result = new Date();
    result.setSeconds(result.getSeconds() + seconds);
    return result;
}

type AppProps = {
    playerName: string;
}
type AppState = {
    phase: string;
    timeRemaining: number;
    alarmTimesRemaining: number[];
    killedBy: string;
    timeToMyImpact?: number;
}

export class App extends React.Component<AppProps, AppState> {
    private updaterId?: number;
    componentWillMount() {
        this.updaterId = setInterval(this.fetchData.bind(this), 1000);
    }
    componentWillUnmount() {
        if (this.updaterId) {
            clearInterval(this.updaterId);
        }
    }
    render() {
        let displayed: any;
        console.log("Rendering:", this.state);
        if (!this.state) {
            return "Loading...";
        }
        if (this.state.killedBy.length > 0) {
            displayed = `Killed by ${this.state.killedBy}`;
        } else {
            switch (this.state.phase) {
                case Phase.PRESTART:
                    displayed = "Game not yet started.";
                    break;
                case Phase.RUNNING:
                    displayed = <div>
                        <LaunchOrConcealButton playerName={this.props.playerName} launched={!!this.state.timeToMyImpact}/>
                        Time remaining: <Timer zeroTime={nowPlus(this.state.timeRemaining)} /> <br />
                        Timers:
                        <ol>
                            {this.state.alarmTimesRemaining.map((d, i) => <li key={i}><Timer zeroTime={nowPlus(d)} /></li>)}
                        </ol>
                    </div>;
                    break;
                case Phase.OVERTIME:
                    displayed = "TODO";
                    break;
                case Phase.ENDED:
                    displayed = `Game over! You're alive! Everyone else is ${this.state.timeToMyImpact ? "dead. Remember? You killed them." : "alive too!"}`;
                    break;
                default:
                    displayed = `Unknown phase: ${this.state.phase}`;
                    break;
            }
        }
        return <div>
            {displayed}
            <pre>{JSON.stringify(this.state, null, 2)}</pre>
        </div>;
    }

    fetchData() {
        jQuery.get({
            url: `/${this.props.playerName}/status`,
            success: dataText => {
                const data = JSON.parse(dataText);
                console.log(data);
                this.setState({
                    phase: data.Phase,
                    timeRemaining: data.TimeRemaining / 1e9,
                    alarmTimesRemaining: (data.AlarmTimesRemaining || []).map((x:number) => x/1e9),
                    killedBy: data.KilledBy,
                    timeToMyImpact: data.TimeToMyImpact || null,
                });
            },
        });
    }
}
