export interface User {
    uid: string,
    email: string,
    displayName: string,
    photoURL: string
}

export interface Session {
    owner: User,
    password: string,
    from: User | null,
    message: any | null
}