export class Constructable {
    type: string = '';
    production: number = 0.0;
    happiness_change: number = 0.0;
    build_cost: number = 0.0;
    build_time: number = 0;
    icon: string = '';
    // production happiness_change population_change build_cost build_time
    stats: Map<string, number[]> = new Map([
        ['Apartment', [0.0, 500.00, 2, 5000, 400000.00, 60]],
        ['Hospital', [1000.00, 5, 1000, 250000.00, 60]],
        ['School', [2000.00, 3, 500, 250000.00, 60]],
        ['Supermarket', [10000.00, 1, 250, 250000.00, 120]],
        ['Barracks', [1000.00, 3, 500, 300000.00, 120]]
    ]);
    icons: Map<string, string> = new Map([
        ['Apartment', '🏢'],
        ['Hospital', '🏥'],
        ['School', '🏫'],
        ['Supermarket', '🏪'],
        ['Barracks', '🎪']
    ]);

    constructor(
        type: string = ''
    ) {
        if (type != '') {
            this.type = type;
            let stats: number[] = this.stats.get(type) as number[];
            this.production = stats[0];
            this.happiness_change = stats[1];
            this.build_cost = stats[3];
            this.build_time = stats[4];
            this.icon = this.icons.get(type) as string;
        }
    }
}
