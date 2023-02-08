# User Stories
### Front-end
As a player, I want to be able to view a page that allows me to create an account so that I can join the game.
As a player, I want to be able to view a page that allows me to provide an already-created username and password so that I can access my progress in the game. 
As a player, I want to view a map of my city that shows the buildings and resources I have amassed while playing the game so that I can see my progress.
As a player, I want to be able to navigate between pages within the game so that I can access both login and game functionality.
### Back-end
As a player, I want to be able to create an account and have my credentials stored in a secure location so that my account is saved.
As a player, I want my data to be linked to my account so that I can have my progress in the game saved when I log back on.
As a player, I want to have a board exclusive to my account so that I can play at my own pace.

# Issues Addressed
- [Initialize HTTP Server in Go](https://github.com/prayujt/city.io/issues/2)
- [Setup connection to MySQL database using .env file](https://github.com/prayujt/city.io/issues/3)
- [Setup HTTP endpoints for creating a user and validating logins](https://github.com/prayujt/city.io/issues/5)
- [Initialize Account Creation Page](https://github.com/prayujt/city.io/issues/6)
- [Initialize Login Page](https://github.com/prayujt/city.io/issues/7)
- [Initialize Game Board and UI](https://github.com/prayujt/city.io/issues/8)
- [Setup navigation between game and login/account](https://github.com/prayujt/city.io/issues/9)
- [Add routing from login page to game page](https://github.com/prayujt/city.io/issues/13)
- [Add linter for Angular code](https://github.com/prayujt/city.io/issues/14)
- [Add use of .env file for Angular API calls](https://github.com/prayujt/city.io/issues/16)
- [Merge .gitignore files](https://github.com/prayujt/city.io/issues/18)
- [Change “npm run start” script to also open the URL](https://github.com/prayujt/city.io/issues/22)

# Plan
Our plan was to set up the login screen and the initial board. We planned to create a registration page that allows users to create an account. This account will obviously belong solely to the user and it will generate a board exclusive to the user. Users will not be able to create an account with a username that is already stored in the database. Additionally, we planned to add an authentication system that verifies user logins. If the user login is valid, then the user will be able to enter the game and view their city. Otherwise, the user will see a popup that tells them that the login is incorrect. 

# What We Accomplished
We were able to make a fully functional login page. It implemented all of the features that we planned to add, in both the front-end and the back-end.

# What Didn’t Work
We were able to start developing the game user interface. However, there are issues with how the game appears based on the dimensions of the screen. In the next sprint, we will be working on standardizing how the game appears across different screens so that the game appears as we want it to. We also had some difficulties adjusting to writing front-end code using the Angular framework. For example, we struggled with setting colors of components, especially ones in which we implemented the Material UI components. 
