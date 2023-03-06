# Progress
We added game functionality and testing in the back-end. Players can now build buildings in their own city. Currently, there are 5 building types: Hospital, School, Supermarket, Barracks, and Apartment. You can see their individual stats (i.e. Production, Happiness, Buildcost, and Build Time). Additionally, you can improve buildings if their upgrades are in the database. For example, you can upgrade a hospital; however, you can not upgrade a school since a level two school doesn’t exist in the database yet. Players can see, in real time, the progress of construction for a building. Furthermore, players can change their city name and view other players' cities. Of course, they are not allowed to edit other players’ cities in any way. A player’s balance will be shown in the top right corner of the screen.
<br />
<br />
A live demo of this game is available [here](http://game.prayujt.com). 
<br />
Our back-end documentation is available in our [`src/README.md`](https://github.com/prayujt/city.io/blob/master/src/README.md) file. 

# Front-end Tests

## Unit Tests
We added unit tests using the `.spec.ts` files for each component of our app. These tests essentially verified that the components were rendered correctly.

## Cypress
We had a series of tests on the login and register components to verify that the login/register system works. Among the test cases included were:
- Valid account registration
- Missing username on register
- Missing password on register
- Confirm password does not match
- Valid account login
- Missing username on login
- Missing password on login
To validate the correct output, I used Cypress’ intercept() method to read server responses to our requests.

# Back-end Tests
We had a test for every HTTP endpoint that we opened on our API, thus giving us greater than a 1:1 function to test ratio. Examples of these tests include:
- Account Creation
- Duplicate Account Creation
- Account Login
- Incorrect Usernames
- Incorrect Passwords
- Validating Correct Sessions
- Validating Invalid Sessions
- City Information Retrieval
- Invalid City Information Retrieval
- Getting Building in Own City
- Getting Building in Not Owned City
- Building Creation
- Invalid Building Creation
- Building Upgrades
<br />
<br />
These tests were done using Golang’s default testing library. 
We set up these tests to run on every pull request using GitHub Actions so that we can see the results before merging a branch in.
