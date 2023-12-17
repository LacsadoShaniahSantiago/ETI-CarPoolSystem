-- CREATE database carpooling_db;
USE carpooling_db;

DROP TABLE IF EXISTS Account;
DROP TABLE IF EXISTS Trip;
DROP TABLE IF EXISTS Enrol;

CREATE TABLE Account (
PassengerID CHAR(3) NOT NULL PRIMARY KEY, 
FirstName VARCHAR(30) NULL,
LastName VARCHAR(30) NULL,
Contact	CHAR(8) NULL,
EmailAddr VARCHAR(50) NULL,
AccountCreated DATETIME NULL,
UserType CHAR(1) NULL,
LicenseNo VARCHAR(10) NULL,
CarPlateNo VARCHAR(10) NULL
);

CREATE TABLE Trip (
TripID CHAR(3) NOT NULL PRIMARY KEY,
PassengerID CHAR(3) NOT NULL,
PickUpAddr VARCHAR(100) NULL,
AltpickUpAddr VARCHAR(60) NULL,
StartTraveltime DATETIME NULL,
DestinationAddr VARCHAR(100) NULL,
PassengerPax INT NULL
);

CREATE TABLE Enrol (
EnrolID CHAR(3) NOT NULL PRIMARY KEY,
TripID CHAR(3) NOT NULL,
PassengerID CHAR(3) NOT NULL,
TripStatus	CHAR(1) NULL
);

INSERT INTO Account	(PassengerID, FirstName, LastName, Contact, EmailAddr, AccountCreated, UserType, LicenseNo, CarPlateNo)
VALUES				('0', "rider","01","12345678", "rider01@email.com", "2023-12-15 08:00:00", "P", null , null),
					('1', "driver","02","87654321", "driver02@email.com", "2023-12-15 08:00:00" ,"D", "B1234567B", "B1234567B");
                    
INSERT INTO Trip	(TripID, PassengerID, PickUpAddr, AltpickUpAddr, StartTravelTime, DestinationAddr, PassengerPax)
VALUES				('0',"1", "123 Drive St 12, SG123123", "456 Drive St 34, SG123456", '2023-08-03 13:00:00', "Causeway Point", 2),
					('1',"1", "321 Drive St 12, SG123321", "654 Drive St 34, SG123654", '2023-09-03 13:00:00', "Jurong Point", 2);
                    
INSERT INTO Enrol (EnrolID, TripID, PassengerID, TripStatus) 
VALUES			  ('0','0','0', 'F'),
				  ('1','0','1', 'F'),
                  ('2','1','0', 'V');

