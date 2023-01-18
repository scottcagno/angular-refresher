import {EventEmitter, Injectable} from '@angular/core';
import {DataService} from "./data.service";

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  isAuthenticated = false;
  authResultEvent = new EventEmitter<boolean>();
  role !:string;

  constructor(private dataService : DataService) { }

  authenticate(username : string, password : string) {
    console.log(`auth.service.ts.authenticate() -> username=${username}, password=${password}`)
    this.dataService.validateUser(username, password).subscribe(
      next => {
        this.setupRole();
        this.isAuthenticated = true;
        this.authResultEvent.emit(true);
      },
      error => {
        this.isAuthenticated = false;
        this.authResultEvent.emit(false);
      }
    );
  }

  setupRole() {
    this.dataService.getRole().subscribe(
      next => {
        console.log(`auth.service.ts.setupRole() -> role=${next.role}`)
        this.role = next.role;
      }
    )
  }

  logout() {
    this.isAuthenticated = false;
  }

}
