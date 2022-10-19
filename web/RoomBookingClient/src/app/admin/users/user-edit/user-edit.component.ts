import {Component, Input, OnInit} from '@angular/core';
import {User} from "../../../model/User";
import {DataService} from "../../../data.service";
import {Router} from "@angular/router";

@Component({
  selector: 'app-user-edit',
  templateUrl: './user-edit.component.html',
  styleUrls: ['./user-edit.component.css']
})
export class UserEditComponent implements OnInit {

  @Input()
  user !:User;
  formUser !:User; // just for the form
  password !:string;
  password2 !:string;
  message !:string;
  nameIsValid = false;
  passwordsAreValid = false;
  passwordsMatch = false;

  constructor(private dataService :DataService, private router :Router) {
  }

  ngOnInit(): void {
    this.formUser = Object.assign({}, this.user);
    this.checkIfNameIsValid();
    this.checkIfPasswordsAreValid();
  }

  onCancel() {
    this.router.navigate(['admin','users'], {queryParams:{id: this.user.id, action:'view'}});
  }

  onSubmit() {
    // check add new or edit existing
    if (this.formUser.id == null) {
      // add new
      this.dataService.addNewUser(this.formUser, this.password).subscribe((user)=>{
        this.router.navigate(['admin','users'], {queryParams:{id: user.id, action: 'view'}});
      });
    } else {
      // edit exsting
      this.dataService.updateUser(this.formUser).subscribe((user)=>{
        this.router.navigate(['admin','users'], {queryParams:{id: user.id, action: 'view'}});
      });
    }
  }

  checkIfNameIsValid() {
    if (this.formUser.name) {
      this.nameIsValid = this.formUser.name.trim().length > 0;
    } else {
      this.nameIsValid = false;
    }
  }

  checkIfPasswordsAreValid() {
    if (this.formUser.id != null) {
      this.passwordsAreValid = true;
      this.passwordsMatch = true;
    } else {
      this.passwordsMatch = this.password === this.password2;
      if (this.password) {
        this.passwordsAreValid = this.password.trim().length > 0;
      } else {
        this.passwordsMatch = false;
      }
    }
  }

}