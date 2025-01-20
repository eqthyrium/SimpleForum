

-- DROP table if exists user;
-- DROP table if exists post;
-- DROP table if exists comment;
-- DROP table if exists category;
-- Drop table if exists post-category;




CREATE TABLE   Users (
                         UserId INTEGER PRIMARY KEY AUTOINCREMENT,
                         Nickname varchar(50) Not null,
                         Email varchar(50) Not null,
                         Password varchar(256),
                         Role varchar(15) not null
);

CREATE TABLE Posts
(     PostId INTEGER PRIMARY KEY AUTOINCREMENT,
      UserId INT NOT NULL,
      Title VARCHAR(255)  NOT NULL,
      Content TEXT NOT NULL,
      Image VARCHAR(255),
      LikeCount INT DEFAULT 0,
      DislikeCount INT DEFAULT 0,
      CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (UserId) REFERENCES user(UserId) ON DELETE CASCADE
);

CREATE TABLE Categories
(
    CategoryId INTEGER PRIMARY KEY AUTOINCREMENT,
    CategoryName VARCHAR(50) NOT NULL
);

CREATE TABLE PostCategories (
    CategoryId INT NOT NULL,
    PostId INT NOT NULL,
    PRIMARY KEY (CategoryId, PostId),
    FOREIGN KEY (CategoryId) REFERENCES Categories(CategoryId) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (PostId) REFERENCES Posts(PostId) ON DELETE CASCADE ON UPDATE CASCADE
                            );

CREATE TABLE Commentaries (
    CommentId INTEGER PRIMARY KEY AUTOINCREMENT,
    PostId INT NOT NULL,
    UserId INT NOT NULL,
    Content TEXT NOT NULL,
    LikeCount INT DEFAULT 0,
    DislikeCount INT DEFAULT 0,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (PostId) REFERENCES  Posts(PostId) ON DELETE CASCADE ,
    FOREIGN KEY (UserId) REFERENCES Users(UserId) ON DELETE CASCADE  );

CREATE TABLE Reactions (
                           UserId INT NOT NULL,
                           PostId INT,
                           CommentId INT,
                           Action CHAR(1) NOT NULL CHECK (Action IN ('L', 'D', 'C')), -- Action: 'L'ike, 'D'islike, 'C'omment
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
    FOREIGN KEY (PostId) REFERENCES Posts(PostId) ON DELETE CASCADE,
    FOREIGN KEY (CommentId) REFERENCES Commentaries(CommentId) ON DELETE CASCADE,
    CONSTRAINT CK_PostOrComment CHECK (PostId IS NOT NULL OR CommentId IS NOT NULL),
    PRIMARY KEY (UserId, PostId, CommentId)
);

CREATE TABLE Reports (
                           UserId INTEGER NOT NULL,
                           PostId INTEGER NOT NULL,
                           FOREIGN KEY (UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
                           FOREIGN KEY (PostId) REFERENCES Posts(PostId) ON DELETE CASCADE
                       );