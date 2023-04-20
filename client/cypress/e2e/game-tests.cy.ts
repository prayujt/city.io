/// <reference types="cypress" />

describe('Game tests', () => {
    beforeEach(() => {
        cy.visit('localhost:4200/login');
        let username: string = 'train';
        let password: string = 'train';
        cy.get("input[placeholder='Username']").type(username);
        cy.get("input[placeholder='Password']").type(password);
        cy.get('button').contains('Log in').click();
        cy.wait(2000);
    });

    it('Build Building', () => {
        cy.get("td[id='2 2']").click();
        cy.get('mat-expansion-panel').contains('Barracks').click();
        cy.wait(1000);
        cy.get('button').each(($ele, index) => {
            if (index == 1) {
                $ele.click();
            }
        });
    });

    it('Train Troops', () => {
        cy.get("a[matTooltip='Train Troops']").click();
        // cy.get('input[matSliderThumb]').focus().type('{rightarrow}');
        // cy.get('a').contains('Train').click();
        // cy.get("a[matTooltip='Train Troops']").click();
        // cy.wait(2000);
    });

    it('Change City Name', () => {
        cy.get("a[matTooltip='Edit City Name']").click();
        cy.get("input[placeholder='City Name']").type('Cypress City');
        cy.get('button').contains('Change Name').click();
        cy.wait(2000);
    });

    it('Scout Button', () => {
        cy.get("a[matTooltip='Scout']").click();
        cy.get('input[aria-label="City"]').type('prayuj');
        cy.get('button').contains('prayuj - monkee city').click();
        cy.wait(2000);
        // cy.get('mat-progress-bar');
    });

    // it("Upgrade Building", () => {
    //     cy.get("td[id='1 1']").click();
    //     cy.get("button").contains("Upgrade").click();
    //     cy.get("mat-progress-bar");
    // })
});
