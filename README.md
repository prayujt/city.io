# city.io
An idle city-building strategy game written in Go and TypeScript using the Angular framework. <br />
<br />
Build your own City! Start as a tiny town and gradually grow to a thriving metropolis. Collect profits and invest in upgrades that grow your city. While you are away, the denizens of your city will generate profit for you. Fill your city with cool businesses, buildings, and whatever else your heart desires. 

# Game Functions
- The game will take place on a grid-based board representing available land, which players can use to construct different kinds of buildings, such as housing and various kinds of businesses.
- Constructing buildings requires the player to spend DobraCoins, which are passively accrued over time based on a variety of factors, including population, residential and commercial buildings, etc. Buildings take time to construct.
- There is a population count that increases or decreases based on certain variables. A higher population count means more DobraCoins! Building housing increases the maximum population count for your city.
- You can spend your DobraCoins to build your army! With this army, you can attack other players to steal their DobraCoins, or fight with the other players for control over the dozens of neutral towns!

# Setup
To run city.io locally: 

For the frontend,
```
$ cd city.io/client
$ npm install
$ npm run start
```
The game will be running on [`localhost:4200`](http://localhost:4200/).

Next, to run the backend,
```
$ cd city.io/src
$ go get
$ go build
$ ./api
```
The api will be running on [`localhost:8000`](http://localhost:8000/).

Alternatively, visit us [`here`](https://cityio.prayujt.com) for our production environment.

# Team
### Back-end
 - Jason Jiang
 - Prayuj Tuli
### Front-end
 - Michael Tu
 - Matthew Mapa
