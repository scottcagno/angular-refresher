import {EventEmitter, Injectable} from '@angular/core';
import {DataService} from "./data.service";

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  isAuthenticated = false;
  authResultEvent = new EventEmitter<boolean>();

  constructor(private dataService : DataService) { }

  authenticate(username : string, password : string) {
    this.dataService.validateUser(username, password).subscribe(
      next => {
        this.isAuthenticated = true;
        this.authResultEvent.emit(true);
      },
      error => {
        this.isAuthenticated = false;
        this.authResultEvent.emit(false);
      }
    );
  }

  logout() {
    this.isAuthenticated = false;
  }

  getRole() :string {
    // if (this.jwtToken == null) {
    //   return ''
    // }
    // // grab the middle, aka payload
    // const encodedPayload = this.jwtToken.split('.')[1];
    // // base64 decode the payload
    // const payload = atob(encodedPayload);
    // // payload is a json string
    // return JSON.parse(payload).role;
    return 'ROLE_ADMIN'; // temp fix
  }
}
