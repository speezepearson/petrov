# Petrov Day Button Game

Developed by Spencer Pearson and William Ehlhardt.

A webapp to play mutually assured destruction between Petrov Day parties, as played at the Seattle/Oxford 2018 Petrov Day parties on September 22, 2018. Features a button, an early warning system, the opportunity to launch a second strike, and, of course, false alarms.

![Screenshot](https://raw.githubusercontent.com/speezepearson/petrov/master/screenshot.png)

# Game setup

Suppose you own a server with a public IP address, `123.45.67.89`, and that you can SSH into `petrov@123.45.67.89`. Then, to get an instance of the game running:

1. **Install Go on the server.** `ssh petrov@123.45.67.89 'sudo apt install golang'` should do it.
2. **Clone this repo.** (On the server or not -- doesn't matter.) `git clone --depth 1 https://github.com/speezepearson/petrov.git` should do it.
3. **Deploy.** Something like `cd petrov; ./ops/launch petrov 123.45.67.89 Seattle:secret-seattle-password Oxford:secret-oxford-password` should do it. This will rsync the code onto `petrov@123.45.67.89`, SSH in, and start the game running, printing out URLs like  `http://123.45.67.89/secret-seattle-password` that you can pass along to the Seattle/Oxford partymasters.

# Player instructions

```
                          TOP SECRET



           MISSILE DEFENSE MANUAL AND STANDING ORDERS

Nuclear warfare is an immensely complicated technical
endeavor. Thankfully, your console reduces it to point-and-click
simplicity.

Your early warning system will detect each launch and immediately
calculate a time to impact. You have until the impact timer reaches 0
to launch a second strike; once the first missiles land, it will be
too late to initiate a response. Launching won't save you from
destruction, but it will at least grant you vengeance.

This satellite-based early warning system is state-of-the-art
technology and totally infallible. Our nation's stance is to launch on
warning, and your orders are to follow that to the letter by launching
immediately in response to any detected missiles. Remember: the
credibility of our nation's second strike, and hence its safety from
nuclear annihilation, rests in your hands.


                            DETAILS

Once you launch, your opponent will immediately be notified, and the
console will show a timer below the launch button reporting the time
until your missile lands.

You may click on the timer to reset your own console to appear as if
you have not launched - rest assured, your missile is still in the
air! Regrettably, your opponent's early warning system will still see
your attack, but hopefully a clean console will help you convince your
opponent of your peaceful intent. Clicking "launch" again will bring
back the time to impact.

If there are missiles in flight at the end of game, the game will
proceed to overtime until all of them land.

You can't unlaunch a missile. What is this, a video game? Treat this
decision with the seriousness it deserves.



                          TOP SECRET

```
