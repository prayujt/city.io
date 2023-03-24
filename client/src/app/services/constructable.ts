interface building {
    buildingType: string;
    buildCost: number;
    buildTime: number;
    buildingProduction: number;
    happinessChange: number;
}

export class Constructable {
    type: string = '';
    production: number = 0.0;
    happiness_change: number = 0.0;
    population_change: number = 0;
    build_cost: number = 0.0;
    build_time: number = 0;
    icon: string = '';

    icons: Map<string, string> = new Map([
        ['Apartment', '🏢'],
        ['Hospital', '🏥'],
        ['School', '🏫'],
        ['Supermarket', '🏪'],
        ['Barracks', '🎪']
    ]);

    constructor(
        building: building
    ) {
        if (building.buildingType != '') {
            this.type = building.buildingType;
            this.production = building.buildingProduction;
            this.happiness_change = building.happinessChange;
            this.build_cost = building.buildCost;
            this.build_time = building.buildTime;
            this.icon = this.icons.get(this.type) as string;
        }
    }
}
