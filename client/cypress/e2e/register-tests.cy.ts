/// <reference types="cypress" />

describe('Register tests', () => {
  beforeEach(() => {
    cy.visit('localhost:4200/login');
    cy.contains("Register").click();
  })

  it('Register missing username', () => {
    cy.get("input[placeholder='Password']").type("password");
    cy.get("input[placeholder='Confirm Password']").type("password");
    cy.get("button").contains("Register").click();
    cy.intercept({
      method: 'POST',
      url: '/login/createAccount'
    }, (req) => {
      req.on('response', (res) => {
        if (res.body['status'] === true)
          throw new Error("should be invalid!");
      })
    })
  })

  it('Register passwords do not match', () => {
    cy.get("input[placeholder='Username']").type("User2");
    cy.get("input[placeholder='Password']").type("password");
    cy.get("input[placeholder='Confirm Password']").type("assword");
    cy.get("button").contains("Register").click();
    cy.intercept({
      method: 'POST',
      url: '/login/createAccount'
    }, (req) => {
      req.on('response', (res) => {
        if (res.body['status'] === true)
          throw new Error("should be invalid!");
      })
    })
  })

  it('Register missing password', () => {
    cy.get("input[placeholder='Username']").type("User3");
    cy.get("button").contains("Register").click();
    cy.intercept({
      method: 'POST',
      url: '/login/createAccount'
    }, (req) => {
      req.on('response', (res) => {
        if (res.body['status'] === true)
          throw new Error("should be invalid!");
      })
    })
  })

  it('Valid Registration and Login', () => {
    let username: string = "User20034"
    let password: string = "password123"
    cy.get("input[placeholder='Username']").type(username);
    cy.get("input[placeholder='Password']").type(password);
    cy.get("input[placeholder='Confirm Password']").type(password);
    cy.get("button").contains("Register").click();
    cy.intercept({
      method: 'POST',
      url: '/login/createAccount'
    }, (req) => {
      req.on('response', (res) => {
        if (res.body['status'] === false)
          throw new Error(res.body['status']);
      })
    })
    cy.wait(1000);
    cy.get("input[placeholder='Username']").type(username);
    cy.get("input[placeholder='Password']").type(password);
    cy.get("button").contains("Log in").click();
  })
})
