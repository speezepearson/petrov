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
    gameEndTime?: Date;
    alarmImpactTimes: Date[];
    killedBy: string;
    myImpactTime?: Date;
    currentTime: Date;
}

export class App extends React.Component<AppProps, AppState> {
    private fetcherId?: number;
    private tickerId?: number;
    constructor(props: AppProps) {
      super(props);
      this.state = {
        phase: undefined,
        gameEndTime: undefined,
        alarmImpactTimes: [],
        killedBy: '',
        myImpactTime: undefined,
        currentTime: new Date(),
      };
    }
    componentWillMount() {
        this.fetcherId = window.setInterval(this.fetchData.bind(this), 1000);
        this.tickerId = window.setInterval(() => this.setState({currentTime: new Date()}), 10);
    }
    componentWillUnmount() {
        if (this.fetcherId) {
            clearInterval(this.fetcherId);
        }
        if (this.tickerId) {
            clearInterval(this.tickerId);
        }
    }
    render() {

        if (this.state.killedBy.length > 0) {
            return `You were killed by ${this.state.killedBy}.`
        }

        switch (this.state.phase) {

            case undefined:
                return 'Loading...';

            case Phase.PRESTART:
                return 'Game not yet started.';

            case Phase.RUNNING:
                const incoming: boolean = (this.state.alarmImpactTimes.length > 0);
                return <div>
                    <div id="time-remaining">
                        {this.state.gameEndTime ? <div><Timer currentTime={this.state.currentTime} zeroTime={this.state.gameEndTime} /> remaining</div> : ''}
                    </div>

                    <div id="top-stuff">
                        {
                            incoming
                            ? [
                                <div id="incoming-label" key="INCOMING">INCOMING</div>,
                                <div id="incoming-timers" key="incoming-timers">
                                    {this.state.alarmImpactTimes.map((d, i) => (
                                        <div className="incoming-timers__timer-wrapper" key={i}>
                                            <div className="incoming-timers__timer">
                                                <Timer currentTime={this.state.currentTime} zeroTime={d} />
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
                            impactTime={this.state.myImpactTime || null}
                            currentTime={this.state.currentTime}
                        />
                    </div>
                </div>;

            case Phase.ENDED:
                return `Game over! You're alive! Everyone else is ${this.state.myImpactTime ? "dead. Remember? You killed them." : "alive too!"}`;

            default:
                return `Unknown phase: ${this.state.phase}`;
        }
    }

    fetchData() {
        jQuery.get({
            url: `/${this.props.playerName}/status`,
            success: dataText => {
                const data = JSON.parse(dataText);
                const now = new Date();
                console.log("received", data, "at", now);
                this.setState({
                    phase: data.Phase,
                    gameEndTime: nowPlus(data.TimeRemaining / 1e9),
                    alarmImpactTimes: (data.AlarmTimesRemaining || []).map((x:number) => nowPlus(x/1e9)),
                    killedBy: data.KilledBy,
                    myImpactTime: data.TimeToMyImpact ? nowPlus(data.TimeToMyImpact / 1e9) : undefined,
                    currentTime: now,
                });
            },
        });
    }
}
