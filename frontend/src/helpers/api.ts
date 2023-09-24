import axios from "axios";
import { UserData } from "../types/user";

const backendPath = "http://localhost:4000"

export function signup(email: string, firstName: string, lastName: string, password: string) {
  return axios.post<UserData>(backendPath + "/api/v1/users/register/", {
    email: email,
    first_name: firstName,
    last_name: lastName,
    password: password,
  });
}

export function activateUser(userId: string, token: string) {
  return axios.put(backendPath + `/api/v1/users/activate/${userId}`, {
    token: token.split(' ').join('')
  });
}

export function signIn(email: string, password: string) {
  return axios.post<UserData>(backendPath + "/api/v1/users/login/", {
    email: email,
    password: password,
  }, {
    withCredentials: true
  });
}

export function getCurrentUser() {
  return axios.get<UserData>(backendPath + "/api/v1/users/current-user", {
    withCredentials: true
  });
}

// TODO: signOut logic on backend side
export function signOut() {
  return console.log("signOut");
}