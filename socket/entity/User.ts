import { Field, Int, ObjectType } from "type-graphql";
import { Entity, PrimaryGeneratedColumn, Column, BaseEntity } from "typeorm";
// import { ObjectType } from "type-graphql";

@ObjectType()
@Entity("users")
export class User extends BaseEntity {
    @Field(() => Int)
    @PrimaryGeneratedColumn()
    id: number

    @Field()
    @Column("text")
    usersName: string

    @Field()
    @Column("text")
    usersRName: string

    @Field()
    @Column("text")
    usersEmail: string

    @Field()
    @Column("text")
    usersPhone: string

    @Column("text")
    usersPassword: string

    @Field()
    @Column("text")
    users2FA: string

    @Field()
    @Column("int")
    usersSuccessLogins: number

    @Field()
    @Column("int")
    usersFailedLogins: number

    @Field()
    @Column("text")
    usersPopularSites: string

    @Field()
    @Column("text")
    usersCurrentChallenge: string

    @Column("int", {default: 0})
    tokenVersion: number;
}