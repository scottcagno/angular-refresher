import { Injectable } from '@angular/core';
import {Layout, LayoutCapacity, Room} from "./model/Room";
import {User} from "./model/User";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  rooms !: Array<Room>;
  users !: Array<User>;

  constructor() {
    // initialize rooms array
    this.rooms = new Array<Room>();
    // initialize users array
    this.users = new Array<User>();
    // populate arrays with initial values
    this.addInitialRooms();
    this.addInitialUsers();
  }

  addInitialRooms() {
    // create room one and layouts for first room
    const room1 = new Room(1, 'First Room', 'First floor');
    const cap1 = new LayoutCapacity(Layout.THEATER, 50);
    const cap2 = new LayoutCapacity(Layout.USHAPE, 20);
    room1.capacities.push(cap1, cap2);
    // create room two and layouts for second room
    const room2 = new Room(2, 'Second Room', 'Third floor');
    const cap3 = new LayoutCapacity(Layout.THEATER, 60);
    room2.capacities.push(cap3);
    // add rooms to rooms array
    this.rooms.push(room1, room2);
  }

  addInitialUsers() {
    // create a few a users
    const user1 = new User(1, "Matt Greencroft");
    const user2 = new User(2, "Laura Croft");
    const user3 = new User(3, "Dick Chesterwood");
    // add users to users array
    this.users.push(user1, user2, user3);
  }
}
