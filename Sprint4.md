# Progress
We have implemented a significant amount of game functionality and balancing, including changes to training and destroying buildings. Players are now able to destroy buildings and receive a small portion of the building cost back. For example, a level three supermarket costs 250,000 dobra coins to build and 750,000 dobra coins to upgrade fully. Selling a level 3 supermarket will return 500,000 dobra coins and also destroy the building. In addition, we have added in-game sounds for building, logging out, and destroying, and barracks now restrict training. Building more barracks increases the rate of training, and upgrading the level of the barracks increases training capacity. Players cannot train without barracks, even with a massive amount of money. Furthermore, players can now view the cities/towns they own and their information, such as city production, troop amount, and population, under the section named "your territory". Visiting conquered towns is now much easier, and players can reinforce their towns with units to defend against enemy attacks. We have also completely fixed marches, as there were some bugs with returning and attacking troops. Additionally, players now have a battle log that allows them to see who they attacked, and more importantly, who attacked them. However, if the player does not win the battle, they will not be able to see the troop count of the other side. Winning a battle allows players to loot the other person and conquer their town.
<br />
<br />
A live demo of this game is available [here](http://cityio.prayujt.com). 
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
A lot of tests had to be fixed this Sprint because of the new changes. Originally, our tests influenced each other. A march test required a training test to function and destroying/upgrading buildings required another test to run first. So, we decided to fix most of the test cases and create separate accounts for each test case. Among the fixed test cases were:
* TestMarch
* TestMarchSuccess
* TestTrainingSuccess
* TestStartMarchAttackPass
* TestTrainingFail
* TestDestroyBuilding
* TestBuildingCreate
* TestBuildingCreateDuplicate
* TestUpgradeBuilding
* TestNameChange
<br />
<br />
These tests were done using Golang’s default testing library. 
We set up these tests to run on every pull request using GitHub Actions so that we can see the results before merging a branch in.


# Running the Game
Our game is publicly accessible at https://cityio.prayujt.com. To run the game locally, follow the instructions in the README.md file after cloning the repository.