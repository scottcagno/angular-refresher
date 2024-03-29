import { Injectable } from '@angular/core';
import {Layout, LayoutCapacity, Room} from "./model/Room";
import {User} from "./model/User";
import {map, Observable, of} from "rxjs";
import {Booking} from "./model/Booking";
import {formatDate} from "@angular/common";
import {environment} from "../environments/environment";
import {HttpClient, HttpHeaders} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  validateUser(username : string, password : string) :Observable<{result:string}> {
    const headers = new HttpHeaders();
    console.log('got creds: ' + username + ', ' + password );
    const creds = btoa(`${username}:${password}`)
    console.log('---> ' + creds)
    headers.append('Authorization','Basic '+creds);
    // go submit the token to the validate method
    return this.http.get<{result:string}>(environment.restUrl + '/api/auth/register', {headers: headers, withCredentials:true});
  }

  /*
   * Room methods
   */
  getRooms() :Observable<Array<Room>> {
    return this.http.get<Array<Room>>(Room.endpoint(), {withCredentials:true}).pipe(
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

  addNewRoom(newRoom: Room) : Observable<Room> {
    return this.http.post<Room>(Room.endpoint(), newRoom);
  }

  updateRoom(room: Room): Observable<Room> {
    return this.http.put<Room>(Room.endpoint(room.id), room, {withCredentials:true});
  }

  deleteRoom(id: number) : Observable<any> {
    return this.http.delete<Room>(Room.endpoint(id))
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

  getRole() :Observable<{ role: string }>{
    const o = this.http.get<{role:string}>(User.getRoleURL(), {withCredentials:true})
    console.log(JSON.stringify(o))
    return o
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
    const fullUser = {id:newUser.id, name:newUser.name, password:password}
    return this.http.post<User>(User.endpoint(), fullUser);
  }

  deleteUser(id: number) : Observable<any> {
    return this.http.delete<User>(User.endpoint(id));
  }

  resetPassword(id: number) :Observable<any> {
    return this.http.get<User>(User.resetPassword(id));
  }


  /*
   * Booking methods
   */
  getBooking(id: number) : Observable<Booking> {
    return this.http.get<Booking>(Booking.endpoint(id))
      .pipe(
        map( data => Booking.fromHttp(data))
      )
  }

  getBookings(date: string) : Observable<Array<Booking>> {
    return this.http.get<Array<Booking>>(Booking.endpointByDate(date))
      .pipe(
        map(data => {
        const bookings = new Array<Booking>();
          for (const booking of data) {
            bookings.push(Booking.fromHttp(booking));
          }
          return bookings as Array<Booking>
    })
    );
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
    return this.http.delete<Booking>(Booking.endpoint(id));
  }


  /*
   * DataService constructor
   */
  constructor(private http: HttpClient) {
    console.log(environment.restUrl);
  }

}
