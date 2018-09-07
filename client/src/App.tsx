import jQuery from 'jquery';
import React from 'react';

import { Timer } from './Timer'
import {LaunchOrConcealButton} from "./LaunchOrConcealButton";
import './App.css';

enum Phase {
    ENDED = "Ended",
    PRESTART = "PreStart",
    RUNNING = "Running",
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
    phase?: string;
    timeRemaining: number;
    alarmTimesRemaining: number[];
    killedBy: string;
    timeToMyImpact?: number;
}

export class App extends React.Component<AppProps, AppState> {
    private updaterId?: number;
    constructor(props: AppProps) {
      super(props);
      this.state = {
        phase: undefined,
        timeRemaining: 0,
        alarmTimesRemaining: [],
        killedBy: '',
        timeToMyImpact: undefined,
      };
    }
    componentWillMount() {
        this.updaterId = setInterval(this.fetchData.bind(this), 1000);
    }
    componentWillUnmount() {
        if (this.updaterId) {
            clearInterval(this.updaterId);
        }
    }
    render() {
        console.log("Rendering:", this.state);

        if (this.state.killedBy.length > 0) {
            return `You were killed by ${this.state.killedBy}.`
        }

        switch (this.state.phase) {

            case undefined:
                return 'Loading...';

            case Phase.PRESTART:
                return 'Game not yet started.';

            case Phase.RUNNING:
                let incoming: boolean = (this.state.alarmTimesRemaining.length > 0);
                return <div>
                    <div id="time-remaining">
                        <Timer zeroTime={nowPlus(this.state.timeRemaining)} /> remaining
                    </div>

                    <div id="top-stuff">
                        {
                            incoming
                            ? [
                                <div id="incoming-label" key="INCOMING">INCOMING</div>,
                                <div id="incoming-timers" key="incoming-timers">
                                    {this.state.alarmTimesRemaining.map((d, i) => (
                                        <div className="incoming-timers__timer-wrapper" key={i}>
                                            <div className="incoming-timers__timer">
                                                <Timer zeroTime={nowPlus(d)} />
                                            </div>
                                        </div>
                                    ))}
                                </div>,
                                <div id="launch-now-label" key="LAUNCH NOW">LAUNCH NOW</div>
                              ]
                            : ''
                        }
                    </div>

                    <div id="bottom-stuff">
                        <LaunchOrConcealButton
                            playerName={this.props.playerName}
                            impactTime={this.state.timeToMyImpact ? nowPlus(this.state.timeToMyImpact) : null}
                        />
                    </div>
                </div>;

            case Phase.ENDED:
                return `Game over! You're alive! Everyone else is ${this.state.timeToMyImpact ? "dead. Remember? You killed them." : "alive too!"}`;

            default:
                return `Unknown phase: ${this.state.phase}`;
        }
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
