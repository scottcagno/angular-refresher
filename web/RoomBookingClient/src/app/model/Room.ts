import {environment} from "../../environments/environment";

export class Room {

  id!: number;
  name!: string;
  location!: string;
  capacities!: Array<LayoutCapacity>;

  constructor(id?: number, name?: string, location?: string) {
    if (id) { this.id = id }
    if (name) { this.name = name }
    if (location) { this.location = location }
    this.capacities = new Array<LayoutCapacity>();
  }

  toString() :string {
    return `${this.name} ${this.location}`
  }

  static fromHttp(data: Room) :Room {
    const newRoom = new Room(data.id, data.name, data.location);
    for (const lc of data.capacities) {
      newRoom.capacities.push(LayoutCapacity.fromHttp(lc));
    }
    return newRoom as Room;
  }

  static endpoint(id ?:number):string {
    if (id) {
      return environment.restUrl + `/api/rooms?id=${id}`
    }
    return environment.restUrl + `/api/rooms`
  }

}

export class LayoutCapacity {
  layout!: Layout;
  capacity!: number;

  constructor(layout?: Layout, capacity?: number) {
    if (layout) { this.layout = layout }
    if (capacity) { this.capacity = capacity }
  }

  toString() :string {
    return `${this.layout}`
  }

  static fromHttp(data: LayoutCapacity) :LayoutCapacity {
    return new LayoutCapacity(data.layout, data.capacity);
  }


}

export enum Layout {
  THEATER = 'Theater',
  USHAPE = 'U-Shape',
  BOARD = 'Board Meeting'
}



// export namespace Layout {
//   export function value(layout: string): typeof Layout[keyof typeof Layout] {
//     const val = Object.keys(Layout).find(key => Layout[key as keyof typeof Layout] == layout);
//     return val as typeof Layout[keyof typeof Layout];
//   }
// }

/*
enum LayoutOptions {
  THEATER = 'Theater',
  USHAPE = 'U-Shape',
  BOARD = 'Board Meeting',
}

export type LayoutOptionKey = keyof typeof LayoutOptions
export type LayoutOptionValue = typeof LayoutOptions[keyof typeof LayoutOptions];
export type LayoutOption = [LayoutOptionKey, LayoutOptionValue];

export const Layouts = new Map<LayoutOptionKey,LayoutOptionValue>([
    LayoutOption(LayoutOptions, LayoutOptions.THEATER),
    LayoutOption(LayoutOptions, LayoutOptions.USHAPE),
    LayoutOption(LayoutOptions, LayoutOptions.BOARD),
  ]
);

function LayoutKey(obj: typeof Layouts, value: LayoutOptionValue) :LayoutOptionKey {
  return Object.keys(obj).find(key => obj[key as keyof typeof Layouts] === value) as LayoutOptionKey;
}

function LayoutOption(options: typeof LayoutOptions, value: LayoutOptionValue) :LayoutOption {
  return [Object.keys(options).find(key => options[key as keyof typeof options] === value) as LayoutOptionKey, value];
}
 */

