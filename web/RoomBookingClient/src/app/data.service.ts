import { Injectable } from '@angular/core';
import {Layout, LayoutCapacity, Room} from "./model/Room";
import {User} from "./model/User";
import {map, Observable, of} from "rxjs";
import {Booking} from "./model/Booking";
import {formatDate} from "@angular/common";
import {environment} from "../environments/environment";
import {HttpClient} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  /*
   * Room methods
   */
  getRooms() :Observable<Array<Room>> {
    return this.http.get<Array<Room>>(Room.endpoint()).pipe(
      map(
        data => {
          const rooms = new Array<Room>();
          for (const room of data) {
            rooms.push(Room.fromHttp(room));
          }
          return rooms;
        }
      )
    );
  }

  getRoom(id: number) :Observable<Room> {
    return this.http.get<Room>(Room.endpoint(id))
      .pipe(
        map(data =>{ return Room.fromHttp(data) })
      );
  }

  addNewRoom(newRoom: Room) : Observable<Room> {
    // @ts-ignore
    return of(null);
  }

  updateRoom(room: Room): Observable<Room> {
    // @ts-ignore
    return of(null);
  }

  deleteRoom(id: number) : Observable<Array<User>> {
    return this.http.get<Array<User>>(User.endpoint())
  }


  /*
   * User methods
   */
  getUsers() :Observable<Array<User>> {
   return this.http.get<Array<User>>(User.endpoint()).pipe(
     map(data => {
       const users = new Array<User>();
       for (const user of data) {
         users.push(User.fromHttp(user));
       }
       return users;
     })
   );
  }

  getUser(id: number) :Observable<User> {
    return this.http.get<User>(User.endpoint(id))
      .pipe(
        map(data =>{ return User.fromHttp(data) })
      );
  }

  updateUser(user :User) :Observable<User> {
    return this.http.put<User>(User.endpoint(user.id), user);
  }

  addNewUser(newUser :User, password :string) :Observable<User> {
    // @ts-ignore
    return of(null);
  }

  deleteUser(id: number) : Observable<any> {
    // @ts-ignore
    return of(null);
  }

  resetUserPassword(id: number) :Observable<any> {
    return of(null);
  }


  /*
   * Booking methods
   */
  getBooking(id: number) : Observable<Booking> {
    // @ts-ignore
    return of(null);
  }

  getBookings(date: string) : Observable<Array<Booking>> {
    // @ts-ignore
    return of(null);
  }

  saveBooking(booking :Booking) :Observable<Booking> {
    // @ts-ignore
    return of(null);
  }

  addBooking(newBooking :Booking) :Observable<Booking> {
    // @ts-ignore
    return of(null);
  }

  deleteBooking(id :number) :Observable<any> {
    // @ts-ignore
    return of(null);
  }


  /*
   * DataService constructor
   */
  constructor(private http: HttpClient) {
    console.log(environment.restUrl);
  }

}
