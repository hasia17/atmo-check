export interface Voivodeship {
  id: string;
  name: string;
  cx: number;
  cy: number;
  islandIds?: string[];
}

export const VOIVODESHIPS: Voivodeship[] = [
  { id: 'dolnoslaskie',        name: 'Dolnośląskie',         cx: 426,  cy: 1282 },
  { id: 'kujawsko-pomorskie',  name: 'Kujawsko-Pomorskie',   cx: 880,  cy: 597  },
  { id: 'lubelskie',           name: 'Lubelskie',            cx: 1781, cy: 1172 },
  { id: 'lubuskie',            name: 'Lubuskie',             cx: 274,  cy: 844  },
  { id: 'lodzkie',             name: 'Łódzkie',              cx: 1057, cy: 1075 },
  { id: 'malopolskie',         name: 'Małopolskie',          cx: 1245, cy: 1666 },
  { id: 'mazowieckie',         name: 'Mazowieckie',          cx: 1423, cy: 865  },
  { id: 'opolskie',            name: 'Opolskie',             cx: 725,  cy: 1417 },
  { id: 'podkarpackie',        name: 'Podkarpackie',         cx: 1675, cy: 1629 },
  { id: 'podlaskie',           name: 'Podlaskie',            cx: 1724, cy: 487  },
  { id: 'pomorskie',           name: 'Pomorskie',            cx: 825,  cy: 224  },
  { id: 'slaskie',             name: 'Śląskie',              cx: 972,  cy: 1531 },
  { id: 'swietokrzyskie',      name: 'Świętokrzyskie',       cx: 1345, cy: 1357 },
  { id: 'warminsko-mazurskie', name: 'Warmińsko-Mazurskie',  cx: 1357, cy: 354  },
  { id: 'wielkopolskie',       name: 'Wielkopolskie',        cx: 667,  cy: 815  },
  {
    id: 'zachodniopomorskie',
    name: 'Zachodniopomorskie',
    cx: 302, cy: 399,
    islandIds: [
      'zachodniopomorskie-island-1',
      'zachodniopomorskie-island-2',
      'zachodniopomorskie-island-3',
      'zachodniopomorskie-island-4',
      'zachodniopomorskie-island-5',
    ],
  },
];
