import { Injectable } from '@angular/core';
import {Layout, LayoutCapacity, Room} from "./model/Room";
import {User} from "./model/User";
import {Observable, of} from "rxjs";
import {Booking} from "./model/Booking";
import {formatDate} from "@angular/common";
import {environment} from "../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  validateUser(username : string, password : string) :Observable<string> {
    return of("ok")
  }

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

  getBooking(id: number) : Observable<Booking> {
    return of(this.bookings.find((b)=>{return b.id === id}) as Booking);
  }

  getBookings(date: string) : Observable<Array<Booking>> {
    return of(this.bookings.filter(b => b.date === date));
  }

  saveBooking(booking :Booking) :Observable<Booking> {
    const existingBooking = this.bookings.find(b => b.id === booking.id) as Booking;
    existingBooking.date = booking.date;
    existingBooking.startTime = booking.endTime;
    existingBooking.endTime = booking.endTime;
    existingBooking.title = booking.title;
    existingBooking.layout = booking.layout;
    existingBooking.room = booking.room;
    existingBooking.user = booking.user;
    existingBooking.participants = booking.participants;
    return of(existingBooking);
  }

  addBooking(newBooking :Booking) :Observable<Booking> {
    let id = 0;
    for (const booking of this.bookings) {
      if(booking.id > id) {
        id = booking.id;
      }
    }
    newBooking.id = id+1;
    this.bookings.push(newBooking);
    return of(newBooking)
  }

  deleteBooking(id :number) :Observable<any> {
    const booking = this.bookings.find((b)=>{return b.id === id}) as Booking;
    this.bookings.splice(this.bookings.indexOf(booking), 1);
    return of(null);
  }

  constructor() {

    console.log(environment.restUrl);


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
        date: formatDate(new Date(), 'yyyy-MM-dd', 'en-us'),
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
        date: formatDate(new Date(), 'yyyy-MM-dd', 'en-us'),
        startTime: '14:00',
        endTime: '15:00',
        participants: 5,
      }
    );
    // add bookings to bookings array
    this.bookings.push(booking1, booking2);
  }
}
