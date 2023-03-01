/// <reference types="cypress" />

describe('Login tests', () => {
    beforeEach(() => {
        cy.visit('localhost:4200/login');
    })
  
    it('Login missing username', () => {
        cy.get("input[placeholder='Password']").type("password");
        cy.get("button").contains("Log in").click();
        cy.intercept({
            method: 'POST',
            url: '/login/createSession'
        }, (req) => {
            req.on('response', (res) => {
            if (res.body['sessionID'] != '')
                throw new Error("should be invalid!");
        })})
    })
  
    it('Login missing password', () => {
        cy.get("input[placeholder='Username']").type("User2");
        cy.get("button").contains("Log in").click();
        cy.intercept({
            method: 'POST',
            url: '/login/createSession'
            }, (req) => {
            req.on('response', (res) => {
                if (res.body['sessionID'] != '')
                    throw new Error("should be invalid!");
        })})
    })

    it('Missing username and password', () => {
        cy.get("button").contains("Log in").click();
        cy.intercept({
            method: 'POST',
            url: '/login/createSession',
            }, (req) => {
            req.on('response', (res) => {
                if (res.body['sessionID'] != '')
                    throw new Error("should be invalid!");
        })})
    })

    it('Login invalid username/password', () => {
        cy.get("input[placeholder='Username']").type("jfioejflajslkfjkaef");
        cy.get("input[placeholder='Password']").type("password")
        cy.get("button").contains("Log in").click();
        cy.intercept({
            method: 'POST',
            url: '/login/createSession'
        }, (req) => {
        req.on('response', (res) => {
            if (res.body['sessionID'] != '')
            throw new Error("should be invalid!");
        })})
    })
  
    it('Valid Login', () => {
        let username: string = "test"
        let password: string = "pass"
        cy.get("input[placeholder='Username']").type(username);
        cy.get("input[placeholder='Password']").type(password);
        cy.get("button").contains("Log in").click();
        cy.intercept({
            method: 'POST',
            url: '/login/createSession'
        }, (req) => {
            req.on('response', (res) => {
            if (res.body['sessionID'] != '')
                throw new Error(res.body['status']);
        })})
    })
  })
  