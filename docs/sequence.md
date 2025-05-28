@startuml
actor User

group Purchase Ticket
autonumber
User -> Server : Request Seat Booking
Server <-> Store : Fetch Sections
Server -> Server : Compute Next Available Seat in sections
Server <-> Store : Get Seat details
Server <-> Store : Update Seat availablity and user
Server <-> Store : Update Section seat availablity
Server -> User : Respond with Booking Confirmation
autonumber stop
end

group Show Receipt
autonumber
User -> Server : Request Receipt with UserId
Server <-> Store : Request User data
Server <-> Store : Fetch Receipts of the user
Server -> Server : Validate Response
Server -> User : Response with User Receipts
autonumber stop
end

group Delete Booking
autonumber
User -> Server : Request Booking Deletion with ReceiptId
Server <-> Store : Request Receipt details
Server -> Server : Validate receipt
Server <-> Store : Reset Seat details to available
Server <-> Store : Increment Section seat availablity
Server -> User: Respond with boolean of operation
autonumber stop
end

@enduml