# Progress
We have enhanced the game's functionality by introducing new features and conducting testing in the backend. Players can now train troops in their cities, but their ability to do so will be limited by their wealth and the capacity of their barracks. Moreover, players will be able to choose whether to train troops in one city or all cities. In addition, players can attack and plunder other players' resources (dobracoin) by sending troops from their own cities to others. The success of such attacks will depend on the number of troops and the level of the barracks. Players can also reinforce their cities with troops to protect them from attacks by other players. Each city will have a designated number of stationed units that will defend it against enemy attacks.
<br />
<br />
A live demo of this game is available [here](http://game.prayujt.com).
<br />
Our back-end documentation is available in our [`src/README.md`](https://github.com/prayujt/city.io/blob/master/src/README.md) file.

## Front-end Progress
On the front-end, we improved the appearance of the game and added additional functionality and security to the game. First, regarding the appearance of the game, we cleaned up the appearance of some of the buttons located on the sidebar to create a cohesive-looking user interface. We also made it so that scouting other cities and towns is done by clicking a scout button rather than clicking on the city name. The sidebar now looks much cleaner and is more user-friendly. For additional functionality, we added three new buttons. First, one major update we made to the game was the ability to train troops and attack other towns. To facilitate this, we added a “Train Troops” button to the sidebar. Upon clicking the train troops button, a dialog box will appear prompting the user to mark the number of troops they want to train using a slider. After doing so and clicking Train, the frontend will send a request to the backend server. Clicking the train troops button again while troops are in the process of being trained will show how many troops are being trained and how much longer it will take. The second new button we added is a “Scout” button. This is similar to last sprint, where users can view other Cities and Towns, but now a new button will appear in the sidebar called Attack!. When the user clicks that button, they can select a city or town to send troops from, then how many troops to send before marching those troops to the other City/Town. This march is viewable in the panel that appears upon clicking the “ViewMarches” button. This button allows users to view marches to and from their city. If no marches are incoming or outgoing, the Marches panel will indicate that there are currently no marches. Otherwise, the marches panel will show a list of marches with information about where they are going, how many troops are being sent, and how much longer the march will take.

## Back-end Progress
To start off, we migrated the system’s authentication method to no longer use our sessionId, and to instead use JWT tokens. We found this to be a lot easier than storing each sessionId and mapping it in our database upon each request, since now we would have the player information available when we decode our JWT. We added a lot of endpoints for army management, such as for training troops, moving troops, getting troop information, getting training information, and getting march information. Most of the back-end work, however, centered upon the march logic resolution, as this took very long to manage. We had to handle all of the different cases and outcomes of battles, and determine the type of resolution. This was accomplished through goroutines and complex database queries that would determine if a movement was an attack, compute troop amounts on both sides, ownership of a tile in which a battle is taking place, and any raiding that would happen. 

# Front-end Tests

## Unit Tests
We added unit tests using the `.spec.ts` files for each component of our app. These tests essentially verified that the components were rendered correctly.

## Cypress
We had a series of tests on the game component to verify that various aspects of the game were functional. Among the test cases included were:
- Troop training
- City name change
- Building construction
- Scouting
</br>
I used Cypress to simulate playing the game and looked for specific elements on the screen to validate that the output was correct.

# Back-end Tests
In this sprint, we did not add too many more endpoints, so we only added a few more tests. The majority of this sprint’s work did not come in the form of adding more HTTP endpoints and functions, but rather in background goroutine threads that would manage the complex battle states. The tests that we did add, though, were:
- City Name Change
- Troop Training Validation
- Invalid Marching
- Validating Troop Movement
<br />
These tests were done using Golang’s default testing library.
We set up these tests to run on every pull request using GitHub Actions so that we can see the results before merging a branch in.


