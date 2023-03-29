import { Constructable } from './constructable';

describe('Constructable', () => {
    it('should create an instance', () => {
        expect(
            new Constructable({
                buildingType: 'test',
                buildCost: 0,
                buildTime: 0,
                buildingProduction: 0,
                happinessChange: 0,
            })
        ).toBeTruthy();
    });
});
