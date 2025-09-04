
export const BRUSH_SIZE_NAMES = ['Single (1 hex)', 'Small (3 hexes)', 'Medium (5 hexes)', 'Large (9 hexes)', 'X-Large (15 hexes)', 'XX-Large (25 hexes)'];

export const AllowedUnitIDs = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 37, 38, 39, 40, 41, 44];
export const CityTerrainIds = [1, 2, 3, 6, 16, 20, 21, 25]; // Base, Hospital, Silo, Mines, City, Tower

// Terrain type names mapping
export const TERRAIN_NAMES: { [key: number]: { name: string, icon: string, color: string } } = {
    0: { name: 'Clear', icon: ' ', color: 'text-black-600 dark:text-green-400' },
    1: { name: 'Land Base', icon: '🌱', color: 'text-green-600 dark:text-green-400' },
    2: { name: 'Naval Base', icon: '🏜️', color: 'text-yellow-600 dark:text-yellow-400' },
    3: { name: 'Airport Base', icon: '🌊', color: 'text-blue-600 dark:text-blue-400' },
    4: { name: 'Dessert', icon: '⛰️', color: 'text-gray-600 dark:text-gray-400' },
    5: { name: 'Grass', icon: '🪨', color: 'text-gray-700 dark:text-gray-300' },
    6: { "name": "Hospital", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    7: { "name": "Mountains", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    8: { "name": "Swamp", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    9: { "name": "Forest", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    10: { "name": "Water (Regular)", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    12: { "name": "Lava", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    14: { "name": "Water (Shallow)", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    15: { "name": "Water (Deep)", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    16: { "name": "Missile Silo", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    17: { "name": "Bridge (Regular)", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    18: { "name": "Bridge (Shallow)", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    19: { "name": "Bridge (Deep)", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    20: { "name": "Mines", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    21: { "name": "City", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    22: { "name": "Road", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    23: { "name": "Water (Rocky)", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    25: { "name": "Guard Tower", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
    26: { "name": "Snow", icon: '⛏️ ', color: 'text-orange-600 dark:text-orange-400' },
};

// Player colors - text colors
export const PLAYER_COLORS: { [key: number]: string } = {
    1: 'text-blue-600 dark:text-red-400',
    2: 'text-blue-600 dark:text-blue-400',
    3: 'text-green-600 dark:text-green-400',
    4: 'text-yellow-600 dark:text-yellow-400',
    5: 'text-orange-600 dark:text-orange-400',
    6: 'text-purple-600 dark:text-purple-400',
    7: 'text-pink-600 dark:text-pink-400',
    8: 'text-cyan-600 dark:text-cyan-400',
    9: 'text-cyan-600 dark:text-cyan-400',
    10: 'text-cyan-600 dark:text-cyan-400',
    11: 'text-cyan-600 dark:text-cyan-400',
    12: 'text-cyan-600 dark:text-cyan-400',
};

// Player background colors for status displays  
export const PLAYER_BG_COLORS: { [key: number]: string } = {
    1: 'bg-sky-900 text-red-800 dark:bg-sky-900 dark:text-red-200',
    2: 'bg-red-100 text-blue-800 dark:bg-red-900 dark:text-blue-200',
    3: 'bg-yellow-100 text-blue-800 dark:bg-yellow-900 dark:text-gray-200',
    4: 'bg-gray-100 text-blue-800 dark:bg-gray-900 dark:text-yellow-200',
    5: 'bg-pink-100 text-blue-800 dark:bg-pink-900 dark:text-pink-200',
    6: 'bg-orange-100 text-blue-800 dark:bg-orange-900 dark:text-orange-200',
    7: 'bg-teal-100 text-blue-800 dark:bg-teal-900 dark:text-teal-200',
    8: 'bg-green-100 text-blue-800 dark:bg-green-900 dark:text-green-200',
    9: 'bg-indigo-100 text-blue-800 dark:bg-indigo-900 dark:text-indigo-200',
    10: 'bg-brown-100 text-brown-800 dark:bg-brown-900 dark:text-brown-200',
    11: 'bg-cyan-100 text-blue-800 dark:bg-cyan-900 dark:text-purple-200',
    12: 'bg-purple-100 text-blue-800 dark:bg-purple-900 dark:text-cyan-200',
};

// Unit type names worldping (basic set)
export const UNIT_NAMES: { [key: number]: { name: string, icon: string } } = {
    1: {name: 'Infantry', icon: ''},
    2: {name: 'Mech', icon: ''},
    3: {name: 'Recon', icon: ''},
    4: {name: 'Tank', icon: ''},
    5: {name: 'Medium Tank', icon: ''},
    6: {name: 'Neo Tank', icon: ''},
    7: {name: 'APC', icon: ''},
    8: {name: 'Artillery', icon: ''},
    9: {name: 'Rocket', icon: ''},
    10: {name: 'Anti-Air', icon: ''},
    11: {name: 'Missile', icon: ''},
    12: {name: 'Fighter', icon: ''},
    13: {name: 'Bomber', icon: ''},
    14: {name: 'B-Copter', icon: ''},
    15: {name: 'T-Copter', icon: ''},
    16: {name: 'Battleship', icon: ''},
    17: {name: 'Cruiser', icon: ''},
    18: {name: 'Lander', icon: ''},
    19: {name: 'Sub', icon: ''},
    20: {name: 'Mech', icon: ''},
    21: {name: 'Missile (Std)', icon: ''},
    22: {name: 'Missile (Nuke)', icon: ''},
    24: {name: 'Sailboat', icon: ''},
    25: {name: 'Artillery (Mega)', icon: ''},
    26: {name: 'Artillery (Quick)', icon: ''},
    27: {name: 'Medic', icon: ''},
    28: {name: 'Stratotanker', icon: ''},
    29: {name: 'Engineer', icon: ''},
    30: {name: 'Goliath RC', icon: ''},
    31: {name: 'Tugboat', icon: ''},
    32: {name: 'Sea Mine', icon: ''},
    33: {name: 'Drone', icon: ''},
    37: {name: 'Cruiser', icon: ''},
    38: {name: 'Missile (Anti Air)', icon: ''},
    39: {name: 'Aircraft Carrier', icon: ''},
    40: {name: 'Miner', icon: ''},
    41: {name: 'Paratrooper', icon: ''},
    44: {name: 'Anti Aircraft (Advanced)', icon: ''},
    // Add more unit types as needed
};
