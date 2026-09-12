-- Created by Redgate Data Modeler (https://datamodeler.redgate-platform.com)
-- Last modification date: 2026-09-03 14:57:27.546

-- tables
-- Table: CastMember
CREATE TABLE CastMember (
    CastMemberId serial  NOT NULL,
    CastMemberName varchar(100)  NOT NULL,
    CastMemberRole char(1)  NOT NULL,
    CastMemberGender char  NOT NULL,
    CONSTRAINT CastMember_pk PRIMARY KEY (CastMemberId)
);

-- Table: Film
CREATE TABLE Film (
    FilmId serial  NOT NULL,
    FilmTitle varchar(200)  NOT NULL,
    FilmReleaseDate date  NOT NULL,
    FilmDuration int  NOT NULL,
    CONSTRAINT Film_pk PRIMARY KEY (FilmId)
);

-- Table: FilmCast
CREATE TABLE FilmCast (
    FIlmId int  NOT NULL,
    CastMemberId int  NOT NULL,
    CONSTRAINT FilmCast_pk PRIMARY KEY (FIlmId,CastMemberId)
);

-- Table: FilmGenre
CREATE TABLE FilmGenre (
    Genre_GenreID int  NOT NULL,
    Film_FilmId int  NOT NULL,
    CONSTRAINT FilmGenre_pk PRIMARY KEY (Genre_GenreID,Film_FilmId)
);

-- Table: FilmUserStatus
CREATE TABLE FilmUserStatus (
    Status_StatusId int  NOT NULL,
    Film_FilmId int  NOT NULL,
    User_UserID int  NOT NULL,
    CONSTRAINT FilmUserStatus_pk PRIMARY KEY (Film_FilmId,User_UserID)
);

-- Table: Genre
CREATE TABLE Genre (
    GenreName varchar(30)  NOT NULL,
    GenreID serial  NOT NULL,
    CONSTRAINT Genre_pk PRIMARY KEY (GenreID)
);

-- Table: Review
CREATE TABLE Review (
    ReviewStars smallint  NOT NULL,
    ReviewDescription text  NOT NULL,
    ReviewUserID int  NOT NULL,
    ReviewFilmId int  NOT NULL,
    CONSTRAINT Review_pk PRIMARY KEY (ReviewFilmId,ReviewUserID)
);

-- Table: Status
CREATE TABLE Status (
    StatusId serial  NOT NULL,
    StatusName varchar(20)  NOT NULL,
    CONSTRAINT Status_pk PRIMARY KEY (StatusId)
);

-- Table: User
CREATE TABLE "User" (
    UserID serial  NOT NULL,
    UserName varchar(120)  NOT NULL,
    UserPasswordHashed varchar(200)  NOT NULL,
    UserBirthday date  NOT NULL,
    UserGenre char  NOT NULL,
    CONSTRAINT User_pk PRIMARY KEY (UserID)
);

-- foreign keys
-- Reference: FK_Cast_CastMember (table: FilmCast)
ALTER TABLE FilmCast ADD CONSTRAINT FK_Cast_CastMember
    FOREIGN KEY (CastMemberId)
    REFERENCES CastMember (CastMemberId)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: FK_Cast_Film (table: FilmCast)
ALTER TABLE FilmCast ADD CONSTRAINT FK_Cast_Film
    FOREIGN KEY (FIlmId)
    REFERENCES Film (FilmId)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: FilmGenre_Film (table: FilmGenre)
ALTER TABLE FilmGenre ADD CONSTRAINT FilmGenre_Film
    FOREIGN KEY (Film_FilmId)
    REFERENCES Film (FilmId)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: FilmGenre_Genre (table: FilmGenre)
ALTER TABLE FilmGenre ADD CONSTRAINT FilmGenre_Genre
    FOREIGN KEY (Genre_GenreID)
    REFERENCES Genre (GenreID)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: FilmUserStatus_Film (table: FilmUserStatus)
ALTER TABLE FilmUserStatus ADD CONSTRAINT FilmUserStatus_Film
    FOREIGN KEY (Film_FilmId)
    REFERENCES Film (FilmId)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: FilmUserStatus_Status (table: FilmUserStatus)
ALTER TABLE FilmUserStatus ADD CONSTRAINT FilmUserStatus_Status
    FOREIGN KEY (Status_StatusId)
    REFERENCES Status (StatusId)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: FilmUserStatus_User (table: FilmUserStatus)
ALTER TABLE FilmUserStatus ADD CONSTRAINT FilmUserStatus_User
    FOREIGN KEY (User_UserID)
    REFERENCES "User" (UserID)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Review_Film (table: Review)
ALTER TABLE Review ADD CONSTRAINT Review_Film
    FOREIGN KEY (ReviewFilmId)
    REFERENCES Film (FilmId)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Review_User (table: Review)
ALTER TABLE Review ADD CONSTRAINT Review_User
    FOREIGN KEY (ReviewUserID)
    REFERENCES "User" (UserID)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- End of file.

