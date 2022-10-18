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
  message !:string;

  constructor(private dataService :DataService, private router :Router) {
  }

  ngOnInit(): void {
    this.formUser = Object.assign({}, this.user);
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

}
