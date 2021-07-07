package service

import (
	"ads-service/domain/model"
	"ads-service/domain/repository"
	"ads-service/domain/service-contracts"
	"ads-service/service/intercomm"
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"time"
)

type campaignService struct {
	repository.CampaignRepository
	repository.InfluencerCampaignRepository
	repository.CampaignUpdateRequestsRepository
	intercomm.AuthClient
	intercomm.UserClient
	intercomm.StoryClient
	intercomm.PostClient
}

func NewCampaignService(r repository.CampaignRepository, ic repository.InfluencerCampaignRepository, curr repository.CampaignUpdateRequestsRepository, ac intercomm.AuthClient, uc intercomm.UserClient, sc intercomm.StoryClient, pc intercomm.PostClient) service_contracts.CampaignService {
	return &campaignService{r , ic,curr, ac, uc, sc, pc}
}


func (c campaignService) GetPostCampaignStatistic(ctx context.Context, bearer string) ([]*model.CampaignStatisticResponse, error) {
	loggedId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return nil, err
	}

	var statResponse []*model.CampaignStatisticResponse

	campaigns, err := c.CampaignRepository.GetAllByOwnerID(ctx, loggedId, "POST")
	if err != nil {
		return nil, err
	}

	influecerCampaigns, err := c.InfluencerCampaignRepository.GetAllByOwnerID(ctx, loggedId, "POST")
	if err != nil {
		return nil, err
	}

	ids := createIdListFromCampaigns(campaigns)
	newIds := createIdListFromInfluencerCampaigns(influecerCampaigns, &ids)

	postMedia, err := c.PostClient.GetPostsFirstMedia(newIds)

	yf,mf,df := time.Now().Date()
	timeef := time.Date(yf,mf,df,0,0,1,0, time.UTC)

	for _, campaign := range campaigns {
		for _, media := range postMedia {
			if media.Id == campaign.ContentId {
				activity := "ACTIVE"
				if campaign.Frequency == "REPEATEDLY" {
					if campaign.DateTo.Before(timeef)  || campaign.DateFrom.After(timeef) {
						activity = "UNACTIVE"
					}
				} else {
					if campaign.ExposeOnceDate != timeef {
						activity = "UNACTIVE"
					}
				}

				statResponse = append(statResponse, &model.CampaignStatisticResponse{
					Id:                       campaign.Id,
					ExposeOnceDate:           campaign.ExposeOnceDate,
					MinDisplaysForRepeatedly: campaign.MinDisplaysForRepeatedly,
					Type:                     campaign.Type,
					Frequency:                campaign.Frequency,
					UserViews: 				  len(campaign.SeenBy),
					WebsiteClicks:            campaign.WebsiteClickCount,
					TargetGroup:              campaign.TargetGroup,
					DateFrom:                 campaign.DateFrom,
					DateTo:                   campaign.DateTo,
					DisplayTime:              campaign.DisplayTime,
					CampaignStatus:           "REGULAR",
					InfluencerUsername:       "",
					Media:                    media.Media,
					Website:                  media.Website,
					Likes:                    media.Likes,
					Dislikes:                 media.Dislikes,
					Comments:                 media.Comments,
					StoryViews:               media.StoryViews,
					Activity:                 model.CampaignStatisticActivity(activity),
					DailyAverage: getDailyAverage(campaign.DailySeenBy),
					InfluencerId: "",
				})
			}
		}
	}

	for _, campaign := range influecerCampaigns {
		for _, media := range postMedia {
			if media.Id == campaign.ContentId {
				statResponse = append(statResponse, &model.CampaignStatisticResponse{
					Id:                       campaign.Id,
					ExposeOnceDate:           time.Now(),
					MinDisplaysForRepeatedly: 0,
					Type:                     campaign.Type,
					Frequency:                "",
					UserViews: 				  len(campaign.SeenBy),
					WebsiteClicks:            campaign.WebsiteClickCount,
					TargetGroup:              model.TargetGroup{
						MinAge: 0,
						MaxAge: 0,
						Gender: "",
					},
					DateFrom:                 time.Now(),
					DateTo:                   time.Now(),
					DisplayTime:              0,
					CampaignStatus:           "INFLUENCER",
					InfluencerUsername:       campaign.Username,
					Media:                    media.Media,
					Website:                  media.Website,
					Likes:                    media.Likes,
					Dislikes:                 media.Dislikes,
					Comments:                 media.Comments,
					StoryViews:               media.StoryViews,
					Activity: "",
					DailyAverage: getDailyAverage(campaign.DailySeenBy),
					InfluencerId: campaign.UserId,
				})
			}
		}
	}

	return statResponse, nil
}

func getDailyAverage(statistic []model.UserGroupStatisticWrapper) float32 {
	if len(statistic) == 0 {
		return 0
	}
	sum := 0
	for _, stat := range statistic {
		sum = sum + len(stat.SeenBy)
	}

	return float32(sum) / float32(len(statistic))
}

func createIdListFromCampaigns(campaigns []*model.Campaign) []string {
	var ids []string
	for _, campaign := range campaigns {
		ids = append(ids, campaign.ContentId)
	}

	return ids
}

func createIdListFromInfluencerCampaigns(campaigns []*model.InfluencerCampaign, ids *[]string) []string {
	for _, campaign := range campaigns {
		*ids = append(*ids, campaign.ContentId)
	}

	return *ids
}

func (c campaignService) GetStoryCampaignStatistic(ctx context.Context, bearer string) ([]*model.CampaignStatisticResponse, error) {
	loggedId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return nil, err
	}

	var statResponse []*model.CampaignStatisticResponse

	campaigns, err := c.CampaignRepository.GetAllByOwnerID(ctx, loggedId, "STORY")
	if err != nil {
		return nil, err
	}

	influecerCampaigns, err := c.InfluencerCampaignRepository.GetAllByOwnerID(ctx, loggedId, "STORY")
	if err != nil {
		return nil, err
	}

	ids := createIdListFromCampaigns(campaigns)
	newIds := createIdListFromInfluencerCampaigns(influecerCampaigns, &ids)

	storyMedia, err := c.StoryClient.GetStoryMedia(newIds)

	yf,mf,df := time.Now().Date()
	timeef := time.Date(yf,mf,df,0,0,1,0, time.UTC)

	for _, campaign := range campaigns {
		for _, media := range storyMedia {
			if media.Id == campaign.ContentId {
				activity := "ACTIVE"
				if campaign.Frequency == "REPEATEDLY" {
					if campaign.DateTo.Before(timeef)  || campaign.DateFrom.After(timeef) {
						activity = "UNACTIVE"
					}
				} else {
					if campaign.ExposeOnceDate != timeef {
						activity = "UNACTIVE"
					}
				}

				statResponse = append(statResponse, &model.CampaignStatisticResponse{
					Id:                       campaign.Id,
					ExposeOnceDate:           campaign.ExposeOnceDate,
					MinDisplaysForRepeatedly: campaign.MinDisplaysForRepeatedly,
					Type:                     campaign.Type,
					Frequency:                campaign.Frequency,
					UserViews:                len(campaign.SeenBy),
					WebsiteClicks:            campaign.WebsiteClickCount,
					TargetGroup:              campaign.TargetGroup,
					DateFrom:                 campaign.DateFrom,
					DateTo:                   campaign.DateTo,
					DisplayTime:              campaign.DisplayTime,
					CampaignStatus:           "REGULAR",
					InfluencerUsername:       "",
					Media:                    media.Media,
					Website:                  media.Website,
					Likes:                    media.Likes,
					Dislikes:                 media.Dislikes,
					Comments:                 media.Comments,
					StoryViews:               media.StoryViews,
					Activity: 				  model.CampaignStatisticActivity(activity),
					DailyAverage: getDailyAverage(campaign.DailySeenBy),
					InfluencerId: "",
				})
			}
		}
	}

	for _, campaign := range influecerCampaigns {
		for _, media := range storyMedia {
			if media.Id == campaign.ContentId {
				statResponse = append(statResponse, &model.CampaignStatisticResponse{
					Id:                       campaign.Id,
					ExposeOnceDate:           time.Now(),
					MinDisplaysForRepeatedly: 0,
					Type:                     campaign.Type,
					Frequency:                "",
					UserViews: 				  len(campaign.SeenBy),
					WebsiteClicks:            campaign.WebsiteClickCount,
					TargetGroup:              model.TargetGroup{
						MinAge: 0,
						MaxAge: 0,
						Gender: "",
					},
					DateFrom:                 time.Now(),
					DateTo:                   time.Now(),
					DisplayTime:              0,
					CampaignStatus:           "INFLUENCER",
					InfluencerUsername:       campaign.Username,
					Media:                    media.Media,
					Website:                  media.Website,
					Likes:                    media.Likes,
					Dislikes:                 media.Dislikes,
					Comments:                 media.Comments,
					StoryViews:               media.StoryViews,
					Activity: "",
					DailyAverage: getDailyAverage(campaign.DailySeenBy),
					InfluencerId: campaign.UserId,
				})
			}
		}
	}

	return statResponse, nil}


func (c campaignService) GetCampaignByPostIdInfulencer(ctx context.Context, contentId string) (*model.CampaignRetreiveRequest, error) {

	campaign, err := c.CampaignRepository.GetFutureByContentIDAndType(ctx, contentId, "POST")
	if err != nil {
		return nil, err
	}

	fmt.Println(campaign)
	return &model.CampaignRetreiveRequest{
		Id:                       campaign.Id,
		MinDisplaysForRepeatedly: campaign.MinDisplaysForRepeatedly,
		Type:                     campaign.Type,
		Frequency:                campaign.Frequency,
		TargetGroup:              campaign.TargetGroup,
		DateFrom:                 campaign.DateFrom,
		DateTo:                   campaign.DateTo,
		DisplayTime:              campaign.DisplayTime,
		ExposeOnceDate:           campaign.ExposeOnceDate,
	}, nil}

func (c campaignService) GetCampaignByPostId(ctx context.Context, bearer string, contentId string) (*model.CampaignRetreiveRequest, error) {
	loggedId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return nil, err
	}

	campaign, err := c.CampaignRepository.GetFutureByContentIDAndType(ctx, contentId, "POST")
	if err != nil {
		return nil, err
	}

	if loggedId != campaign.OwnerId {
		return nil, errors.New("unauthorized to campaign access")
	}

	return &model.CampaignRetreiveRequest{
		Id:                       campaign.Id,
		MinDisplaysForRepeatedly: campaign.MinDisplaysForRepeatedly,
		Type:                     campaign.Type,
		Frequency:                campaign.Frequency,
		TargetGroup:              campaign.TargetGroup,
		DateFrom:                 campaign.DateFrom,
		DateTo:                   campaign.DateTo,
		DisplayTime:              campaign.DisplayTime,
		ExposeOnceDate:           campaign.ExposeOnceDate,
	}, nil}

func (c campaignService) GetCampaignByStoryId(ctx context.Context, bearer string, contentId string) (*model.CampaignRetreiveRequest, error) {
	loggedId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return nil, err
	}

	campaign, err := c.CampaignRepository.GetFutureByContentIDAndType(ctx, contentId, "STORY")
	if err != nil {
		return nil, err
	}

	if loggedId != campaign.OwnerId {
		return nil, errors.New("unauthorized to campaign access")
	}

	return &model.CampaignRetreiveRequest{
		Id:                       campaign.Id,
		MinDisplaysForRepeatedly: campaign.MinDisplaysForRepeatedly,
		Type:                     campaign.Type,
		Frequency:                campaign.Frequency,
		TargetGroup:              campaign.TargetGroup,
		DateFrom:                 campaign.DateFrom,
		DateTo:                   campaign.DateTo,
		DisplayTime:              campaign.DisplayTime,
	}, nil
}

func (c campaignService) DeleteCampaignByPostId(ctx context.Context, bearer string, contentId string) error {
	campaign, err := c.CampaignRepository.GetByContentIDAndType(ctx, contentId, "POST")
	if err != nil {
		return err
	}

	return c.CampaignRepository.DeleteByID(ctx, campaign.Id)
}

func (c campaignService) DeleteCampaignByStoryId(ctx context.Context, bearer string, contentId string) error {
	campaign, err := c.CampaignRepository.GetByContentIDAndType(ctx, contentId, "STORY")
	if err != nil {
		return err
	}

	return c.CampaignRepository.DeleteByID(ctx, campaign.Id)
}

func (c campaignService) ClickOnStoryCampaignWebsite(ctx context.Context, contentId string) error {
	campaign, err := c.CampaignRepository.GetByContentIDAndType(ctx, contentId, "STORY")
	if err != nil {
		return err
	}

	campaign.WebsiteClickCount = campaign.WebsiteClickCount + 1
	_, err = c.CampaignRepository.Update(ctx, campaign)
	if err != nil {
		return err
	}

	return nil
}

func (c campaignService) ClickOnPostCampaignWebsite(ctx context.Context, contentId string) error {
	campaign, err := c.CampaignRepository.GetByContentIDAndType(ctx, contentId, "POST")
	if err != nil {
		return err
	}

	campaign.WebsiteClickCount = campaign.WebsiteClickCount + 1
	_, err = c.CampaignRepository.Update(ctx, campaign)
	if err != nil {
		return err
	}

	return nil
}

func (c campaignService) GetCampaignById(ctx context.Context, bearer string, campaignId string) (*model.Campaign, error) {
	loggedId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return nil, err
	}

	campaign, err := c.CampaignRepository.GetByID(ctx, campaignId)
	if err != nil {
		return nil, err
	}

	if loggedId != campaign.OwnerId {
		return nil, errors.New("unauthorized to campaign access")
	}

	return campaign, nil
}

func (c campaignService) GetAllActiveAgentsPostCampaigns(ctx context.Context, bearer string) ([]string, error) {
	loggedId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return nil, err
	}

	campaigns, err := c.CampaignRepository.GetAllFutureByOwnerIDAndType(ctx, loggedId, "POST")
	if err != nil {
		return nil, err
	}

	return getContentIdsFromCampaigns(campaigns), nil
}

func getContentIdsFromCampaigns(campaigns []*model.Campaign) []string {
	var retVal []string
	for _, campaign := range campaigns {
		retVal = append(retVal, campaign.ContentId)
	}
	if retVal == nil {
		return []string{}
	}
	return retVal
}

func (c campaignService) GetAllActiveAgentsStoryCampaigns(ctx context.Context, bearer string) ([]string, error) {
	loggedId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return nil, err
	}

	campaigns, err := c.CampaignRepository.GetAllFutureByOwnerIDAndType(ctx, loggedId, "STORY")
	if err != nil {
		return nil, err
	}

	return getContentIdsFromCampaigns(campaigns), nil
}

func (c campaignService) CreateCampaign(ctx context.Context, bearer string, campaignRequest *model.CampaignRequest) (string, error) {
	loggedId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return "", err
	}

	campaign, err := model.NewCampaign(campaignRequest, loggedId)
	if err != nil {
		return "", err
	}

	_, err = c.CampaignRepository.Create(ctx, campaign)
	if err != nil {
		return "", err
	}

	return campaign.Id, nil
}

func (c campaignService) GetUnseenPostIdsCampaignsForUser(ctx context.Context, bearer string, count int) ([]string, error) {
	targetUser, err := c.UserClient.GetLoggedUserTargetGroup(bearer)
	if err != nil {
		return []string{}, err
	}

	suggestions, err := c.CampaignRepository.GetUnseenContentIdsCampaignsForUser(ctx, targetUser, "POST", count)
	if err != nil || suggestions == nil{
		return []string{}, err
	}

	var retVal []string
	i := 0

	y,m,d := time.Now().Date()
	today := time.Date(y,m,d,0,0,0,0, time.UTC)

	for _, suggestion := range suggestions {
		if count <= i {
			break
		}
		if !isSeenByUser(suggestion.SeenBy, targetUser.Id) {
			suggestion.SeenBy = append(suggestion.SeenBy, targetUser.Id)
		}
		if !isSeenByUserToday(suggestion.DailySeenBy, targetUser.Id, today){
			if hasDay(suggestion.DailySeenBy, today) {
				idx := 0
				for i, seen := range suggestion.DailySeenBy {
					if seen.Date == today {
						idx = i
						break
					}
				}
				suggestion.DailySeenBy[idx].SeenBy = append(suggestion.DailySeenBy[idx].SeenBy, model.UserGroupStatistic{
					Id:  targetUser.Id,
					Age: targetUser.Age,
				})
			} else {
				suggestion.DailySeenBy = append(suggestion.DailySeenBy, model.UserGroupStatisticWrapper{
					Date:   today,
					SeenBy: []model.UserGroupStatistic{{
						Id:  targetUser.Id,
						Age: targetUser.Age,
					}},
				})
			}
		}

		c.CampaignRepository.Update(ctx, suggestion)
		retVal = append(retVal, suggestion.ContentId)
		i = i + 1
	}

	return retVal, nil
}

func isSeenByUser(seenBy []string, userId string) bool {
	for _, seen := range seenBy{
		if seen == userId {
			return true
		}
	}
	return false
}

func hasDay(seenBy []model.UserGroupStatisticWrapper, today time.Time) bool {
	for _, seen := range seenBy {
		fmt.Println(seen.Date)
		fmt.Println(today)
		if seen.Date == today {
			return true
		}
	}

	return false
}

func isSeenByUserToday(seenBy []model.UserGroupStatisticWrapper, userId string, today time.Time) bool {

	if !hasDay(seenBy, today) {
		return false
	}

	for _, seen := range seenBy {
		if seen.Date == today {
			for _, users := range seen.SeenBy {
				if users.Id == userId {
					return true
				}
			}
		}
	}

	return false
}

func (c campaignService) GetUnseenStoryIdsCampaignsForUser(ctx context.Context, bearer string, count int) ([]string, error){
	targetUser, err := c.UserClient.GetLoggedUserTargetGroup(bearer)
	if err != nil {
		return []string{}, err
	}

	suggestions, err := c.CampaignRepository.GetUnseenContentIdsCampaignsForUser(ctx, targetUser, "STORY", count)
	if err != nil || suggestions == nil {
		return []string{}, err
	}

	var retVal []string
	i := 0

	y,m,d := time.Now().Date()
	today := time.Date(y,m,d,0,0,0,0, time.UTC)

	for _, suggestion := range suggestions {
		if count <= i {
			break
		}
		if !isSeenByUser(suggestion.SeenBy, targetUser.Id) {
			suggestion.SeenBy = append(suggestion.SeenBy, targetUser.Id)
		}
		if !isSeenByUserToday(suggestion.DailySeenBy, targetUser.Id, today){
			if hasDay(suggestion.DailySeenBy, today) {
				idx := 0
				for i, seen := range suggestion.DailySeenBy {
					if seen.Date == today {
						idx = i
						break
					}
				}
				suggestion.DailySeenBy[idx].SeenBy = append(suggestion.DailySeenBy[idx].SeenBy, model.UserGroupStatistic{
					Id:  targetUser.Id,
					Age: targetUser.Age,
				})
			} else {
				suggestion.DailySeenBy = append(suggestion.DailySeenBy, model.UserGroupStatisticWrapper{
					Date:   today,
					SeenBy: []model.UserGroupStatistic{{
						Id:  targetUser.Id,
						Age: targetUser.Age,
					}},
				})
			}
		}

		c.CampaignRepository.Update(ctx, suggestion)
		retVal = append(retVal, suggestion.ContentId)
		i = i + 1
	}

	return retVal, nil
}


func (c campaignService) CreateInfluencerCampaign(ctx context.Context, bearer string, campaignRequest *model.InfluencerCampaignCreateRequest) (string, error) {
	logged, err := c.UserClient.GetLoggedUserInfo(bearer)
	if err != nil {
		return "", err
	}
	fmt.Println("USO2")

	//mozda posle nece trebati provera
	//_, err = c.CampaignRepository.GetByID(ctx, campaignRequest.ParentCampaignId)
	//if err != nil {
	//	return "", errors.New("invalid parent campaign id")
	//}

	campaign, err := model.NewInfluencerCampaign(campaignRequest, logged.Id, logged.Username)
	if err != nil {
		return "", err
	}

	if err = validator.New().Struct(campaign); err!= nil {
		return "", err
	}

	_, err = c.InfluencerCampaignRepository.Create(ctx, campaign)
	if err != nil {
		return "", err
	}

	return campaign.Id, nil
}


func (c campaignService) CreateCampaignForInfluencer(ctx context.Context, campaignRequest *model.InfluencerCampaignProductCreateRequest) (string, error) {


	campaign, err := model.NewInfluencerCampaignProduct(campaignRequest)
	if err != nil {
		return "", err
	}

	if err = validator.New().Struct(campaign); err!= nil {
		return "", err
	}

	_, err = c.InfluencerCampaignRepository.Create(ctx, campaign)
	if err != nil {
		return "", err
	}

	return campaign.Id, nil
}

func (c campaignService) UpdateCampaignRequest(ctx context.Context, bearer string, campaignRequest *model.CampaignUpdateRequestDTO) (string, error) {
	loggedId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return "", err
	}

	campaignReq, err := c.CampaignUpdateRequestsRepository.GetPendingByCampaignId(ctx, campaignRequest.CampaignId)
	if err != nil {
		return "", errors.New("database connection problem")
	}

	if campaignReq != nil {
		return "", errors.New("there is pending request for campaign update")
	}

	campaign, err := c.CampaignRepository.GetByID(ctx, campaignRequest.CampaignId)
	if err != nil {
		return "", errors.New("invalid campaign id")
	}

	if loggedId != campaign.OwnerId {
		return "", errors.New("unauthorized to change campaign")
	}

	if campaign.Frequency != "REPEATEDLY" {
		return "", errors.New("cannot edit campaign that lasts only one day")
	}

	campaignUpdateRequest, err := model.NewCampaignUpdateRequest(campaignRequest)
	if err != nil {
		return "", err
	}

	_, err = c.CampaignUpdateRequestsRepository.Create(ctx, campaignUpdateRequest)
	if err != nil {
		return "", err
	}

	return campaignUpdateRequest.Id, nil
}

func (c campaignService) UpdateCampaignVisitor(ctx context.Context, bearer string, id string, campaignType string) error {
	logged, err := c.UserClient.GetLoggedUserTargetGroup(bearer)
	if err != nil {
		return err
	}

	campaignReq, err := c.CampaignRepository.GetByContentIDAndType(ctx, id, campaignType)
	if err != nil {
		return errors.New("database connection problem")
	}

	y,m,d := time.Now().Date()
	today := time.Date(y,m,d,0,0,0,0, time.UTC)

	if !isSeenByUser(campaignReq.SeenBy, logged.Id) {
		campaignReq.SeenBy = append(campaignReq.SeenBy, logged.Id)
	}

	if !isSeenByUserToday(campaignReq.DailySeenBy, logged.Id, today){
		if hasDay(campaignReq.DailySeenBy, today) {
			idx := 0
			for i, seen := range campaignReq.DailySeenBy {
				if seen.Date == today {
					idx = i
					break
				}
			}
			campaignReq.DailySeenBy[idx].SeenBy = append(campaignReq.DailySeenBy[idx].SeenBy, model.UserGroupStatistic{
				Id:  logged.Id,
				Age: logged.Age,
			})
		} else {
			campaignReq.DailySeenBy = append(campaignReq.DailySeenBy, model.UserGroupStatisticWrapper{
				Date:   today,
				SeenBy: []model.UserGroupStatistic{{
					Id:  logged.Id,
					Age: logged.Age,
				}},
			})
		}
	}

	_, err = c.CampaignRepository.Update(ctx, campaignReq)
	if err != nil {
		return errors.New("database connection problem")
	}

	return nil
}