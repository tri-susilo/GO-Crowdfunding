package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
	FindByID(ID int) (Campaign, error)
	AddCampaign(campaign Campaign) (Campaign, error)
	UpdateCampaign(campaign Campaign) (Campaign, error)
	UploadImage(campaignImage CampaignImage) (CampaignImage, error)
	MarkAllImagesNonPrimary(campaignID int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// repository function to search all campaign
func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// repository function to find campaign by user id
func (r *repository) FindByUserID(UserID int) ([]Campaign, error) {
	var campaigns []Campaign
	// "campaign_images.is_primary = 1"
	err := r.db.Where("user_id=?", UserID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// repository function to Search Campaign by ID
func (r *repository) FindByID(ID int) (Campaign, error) {
	var campaign Campaign
	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", ID).Find(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

// repository function to Create Campaign
func (r *repository) AddCampaign(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

// repository function to Update Campaign
func (r *repository) UpdateCampaign(campaign Campaign) (Campaign, error) {
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

// repository function to upload Image
func (r *repository) UploadImage(campaignImage CampaignImage) (CampaignImage, error) {
	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}
	return campaignImage, nil

}

func (r *repository) MarkAllImagesNonPrimary(campaignID int) (bool, error) {

	//Update campaign_image set is_primary = false where campaign_id = campaign_id

	err := r.db.Model(&CampaignImage{}).Where("campaign_id", campaignID).Update("is_primary", false).Error

	if err != nil {
		return false, err
	}
	return true, nil
}
