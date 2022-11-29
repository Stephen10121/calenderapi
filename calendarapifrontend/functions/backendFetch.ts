import { GroupsType } from "../src/components/Groups";
import { PendingGroupsType } from "../src/components/PendingGroups";
import { POST_SERVER } from "./variables";

export interface FetchGroupsResponse {
    error?: string
    data?: {
        groups: GroupsType[]
        pendingGroups: PendingGroupsType[]
    }
}

export interface LoginResponse {
    error: string
    message?: string
    data?: {
        email: string
        name: string
        token: string
    }
}

export interface ValidateResponse {
    error: boolean
    data?: {
        email: string
        name: string
    }
}

export async function fetchGroups(token: string): Promise<FetchGroupsResponse> {
    try {
      const groups = await fetch(`${POST_SERVER}/myGroups`, {
        method: "GET",
        headers: {
          "Authorization": `Bearer ${token}`
        },
        credentials: "omit"
      })
      const groupsJson = await groups.json();
      if (groupsJson.error) {
          return {error: groupsJson.error}
      }
      return { data: { groups: groupsJson.groups, pendingGroups: groupsJson.pendingGroups} }
    } catch (err) {
      console.error(err);
      return {error: "Failed to fetch Groups."};
    }
}

export async function validate(token: string): Promise<ValidateResponse> {
    try {
      const groups = await fetch(`${POST_SERVER}/validate`, {
        method: "GET",
        headers: {
          "Authorization": `Bearer ${token}`
        },
        credentials: "omit"
      })
      const groupsJson = await groups.json();
      if (groupsJson.error) {
        return {error: true}
      }
      return {error: false, data: groupsJson}
    } catch (err) {
      console.error(err);
      return {error: true};
    }
}

export async function login(email: string, password: string): Promise<LoginResponse> {
    try {
      const groups = await fetch(`${POST_SERVER}/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "omit",
        body: JSON.stringify({
            "email": email,
            "password": password
        })
      })
      const groupsJson = await groups.json();
      if (groupsJson.error) {
        return {error: groupsJson.error}
      }
      return {error: "", data: groupsJson}
    } catch (err) {
      console.error(err);
      return {error: "Error Logging in."};
    }
}