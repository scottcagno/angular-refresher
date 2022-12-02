import {Layout, Room} from "./Room";
import {User} from "./User";
import {DataService} from "../data.service";
import {Router} from "@angular/router";
import {environment} from "../../environments/environment";

export class Booking {
  id !: number;
  room !: Room;
  user !: User;
  layout !: Layout;
  title !: string;
  date !: string;
  startTime !: string;
  endTime !: string;
  participants !: number;

  constructor(data:{id?:number, room?:Room, user?:User, layout?:Layout, title?:string, date?:string,
    startTime?:string, endTime?:string, participants?:number}) {
    if (data.id) { this.id = data.id as number }
    if (data.room) { this.room = data.room }
    if (data.user) { this.user = data.user }
    if (data.layout) { this.layout = data.layout }
    if (data.title) { this.title = data.title }
    if (data.date) { this.date = data.date }
    if (data.startTime) { this.startTime = data.startTime }
    if (data.endTime) { this.endTime = data.endTime }
    if (data.participants) { this.participants = data.participants }
  }

  static fromHttp(b: Booking) :Booking {
    return new Booking({
      id: b.id as number,
      room : Room.fromHttp(b.room) as Room,
      user : User.fromHttp(b.user) as User,
      layout: b.layout,
      title: b.title,
      date : b.date,
      startTime: b.startTime,
      endTime: b.endTime,
      participants: b.participants,
    }) as Booking;
  }

  onCancel() {}

  onSubmit() {}

  // constructor(
  //   id?:number,
  //   room?:Room,
  //   user?:User,
  //   layout?:Layout,
  //   title?:string,
  //   date?:string,
  //   startTime?:string,
  //   endTime?:string,
  //   participants?:number) {
  //   if (id) {this.id = id}
  //   if (room) {this.room = room}
  //   if (user) {this.user = user}
  //   if (layout) {this.layout = layout}
  //   if (date) {this.date = date}
  //   if (startTime) {this.startTime = startTime}
  //   if (endTime) {this.endTime = endTime}
  //   if (participants) {this.participants = participants}
  // }

  static endpoint(id ?:number):string {
    if (id) {
      return environment.restUrl + `/api/bookings?id=${id}`
    }
    return environment.restUrl + `/api/bookings`
  }

  static endpointByDate(date ?:string):string {
    if (date) {
      return environment.restUrl + `/api/bookings?date=${date}`
    }
    return environment.restUrl + `/api/bookings`
  }

  getDateAsDate() {
    return new Date(this.date);
  }
}
