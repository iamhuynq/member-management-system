package setting

// config constant
const (
	/*
		User role: normal user, admin
	*/
	UserRoleType  = 1
	AdminRoleType = 2

	/*
		Limit character for name
	*/
	NameMaxLength = 40

	/*
		Limit character for password
	*/
	PasswordMinLength = 6

	/*
		Member gender: male, female, other
	*/
	MaleGenderType   = 1
	FemaleGenderType = 2
	OtherGenderType  = 3

	/*
		Member status in company: In(Enrolling), Out(Quited)
	*/
	MemberInStatus  = 0
	MemberOutStatus = 1

	/*
		Limit character for account
	*/
	LoginIDMaxLength    = 100
	SNSAccountMaxLength = 100

	/*
		Limit character for additional information
	*/
	CommentMaxLength     = 1000
	DescriptionMaxLength = 1000

	/*
		Limit character for color code (HEX code)
	*/
	ColorCodeMaxLength = 11

	/*
		Limit department in company
	*/
	LeastDepartment = 1
	MaxDepartment   = 20

	/*
		Seat information status: active, inactive, all
	*/
	SeatMasterInactiveStatus = 0
	SeatMasterActiveStatus   = 1
	SeatMasterAllStatus      = 2
	/*
		Position in team: normal member, leader
	*/
	MemberPosition = 1
	LeaderPosition = 2

	/*
		Limit size for upload image
	*/
	FileMaxSize = 3 << 20

	/*
		Memmory size for buffer
	*/
	ReadBufferSize  = 1024
	WriteBufferSize = 1024

	/*
		Sidebar actived tab (for template)
	*/
	AdminHomeTab = 0
	MembersTab   = 1
	TeamsTab     = 2
	CompaniesTab = 3
	SeatsTab     = 4

	/*
		Move type: move on team, move on seat
	*/
	MoveOnTeam = 0
	MoveOnSeat = 1

	/*
		Check user delete or upload avatar
	*/
	DeleteAvatar = "0"
	UploadAvatar = "1"

	/*
		Default template for normal user, admin
	*/
	AdminTemplate = "template/layouts/admin_template.tmpl"
	UserTemplate  = "template/layouts/user_template.tmpl"

	/*
		PATH to default avatar for member, team
	*/
	DefaultLocalMemberAvatar = "webroot/img/avatar_empty.png"
	DefaultLocalTeamAvatar   = "webroot/img/team_default.png"

	/*
		PATH to image folder
	*/
	ImageBaseURL = "webroot/img"

	/*
		Constant for AWS
	*/
	S3BucketURL           = "https://vista-icons.s3-ap-northeast-1.amazonaws.com"
	BucketName            = "vista-icons"
	DefaultS3MemberAvatar = "https://vista-icons.s3-ap-northeast-1.amazonaws.com/members/default.png"
	DefaultS3TeamAvatar   = "https://vista-icons.s3-ap-northeast-1.amazonaws.com/teams/default.png"
	S3MemberFolder        = "/members/"
	S3TeamFolder          = "/teams/"
	MemberFolderType      = 1
	TeamFolderType        = 2
	Region                = "ap-northeast-1"
)
