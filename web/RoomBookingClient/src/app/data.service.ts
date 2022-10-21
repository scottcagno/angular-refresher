import { Injectable } from '@angular/core';
import {Layout, LayoutCapacity, Room} from "./model/Room";
import {User} from "./model/User";
import {Observable, of} from "rxjs";
import {Booking} from "./model/Booking";
import {formatDate} from "@angular/common";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  private rooms !: Array<Room>;
  private users !: Array<User>;
  private bookings !: Array<Booking>;

  getRooms() :Observable<Array<Room>> {
    return of(this.rooms);
  }

  addNewRoom(newRoom: Room) : Observable<Room> {
    let id = 0;
    for (const room of this.rooms) {
      if (room.id > id) {
        id = room.id;
      }
    }
    newRoom.id = id+1;
    this.rooms.push(newRoom);
    return of(newRoom);
  }

  updateRoom(room: Room): Observable<Room> {
    const originalRoom = this.rooms.find(r => r.id === room.id) as Room;
    originalRoom.name = room.name;
    originalRoom.location = room.location;
    originalRoom.capacities = room.capacities;
    return of(originalRoom);
  }

  deleteRoom(id: number) : Observable<any> {
    const room = this.rooms.find(r => r.id === id) as Room;
    this.rooms.splice(this.rooms.indexOf(room), 1);
    return of(null);
  }

  getUsers() :Observable<Array<User>> {
    return of(this.users);
  }

  updateUser(user :User) :Observable<User> {
    const originalUser = this.users.find(u => u.id === user.id ) as User;
    originalUser.name = user.name;
    return of(originalUser);
  }

  addNewUser(newUser :User, password :string) :Observable<User> {
    let id = 0;
    for (const user of this.users) {
      if (user.id > id) {
        id = user.id;
      }
    }
    newUser.id = id+1;
    this.users.push(newUser);
    return of(newUser);
  }

  deleteUser(id: number) : Observable<any> {
    const user = this.users.find(u => u.id === id) as User;
    this.users.splice(this.users.indexOf(user), 1);
    return of(null);
  }

  resetUserPassword(id: number) :Observable<any> {
    return of(null);
  }

  getBookings() : Observable<Array<Booking>> {
    return of(this.bookings);
  }

  constructor() {
    // initialize rooms array
    this.rooms = new Array<Room>();
    // initialize users array
    this.users = new Array<User>();
    // initialize bookings array
    this.bookings = new Array<Booking>();
    // populate arrays with initial values
    this.addInitialRooms();
    this.addInitialUsers();
    this.addInitialBookings();
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

  addInitialBookings() {
    // create a few bookings
    const booking1 = new Booking(
      {
        id: 1,
        room: this.rooms[0],
        user: this.users[0],
        layout: Layout.THEATER,
        title: 'Example meeting',
        date: formatDate(new Date(), 'medium', 'en-us'),
        startTime: '11:30',
        endTime: '12:30',
        participants: 12,
      });
    const booking2 = new Booking(
      {
        id: 2,
        room: this.rooms[1],
        user: this.users[1],
        layout: Layout.USHAPE,
        title: 'Another meeting',
        date: formatDate(new Date(), 'medium', 'en-us'),
        startTime: '14:00',
        endTime: '15:00',
        participants: 5,
      }
    );
    // add bookings to bookings array
    this.bookings.push(booking1, booking2);
  }
}
