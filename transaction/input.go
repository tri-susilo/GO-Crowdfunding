package transaction

import "crowdfunding/user"

type GetCampaignTransactionsinput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
