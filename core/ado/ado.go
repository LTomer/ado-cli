package ado

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strconv"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
)

type Ado struct {
	ctx        context.Context
	connection *azuredevops.Connection
}

func NewClient(url string, pat string) *Ado {
	a := &Ado{
		ctx:        context.Background(),
		connection: azuredevops.NewPatConnection(url, pat),
	}

	return a
}

func (p *Ado) GetPrivateQueries(projectName string) {

	project, err := p.getProjectID(projectName)
	if err != nil {
		log.Fatalf("Error getting project ID: %v", err)
	}

	// List private queries
	privateQueries, err := p.listPrivateQueries(project.Id.String())
	if err != nil {
		log.Fatalf("Error listing private queries: %v", err)
	}

	// Print the private queries
	for _, query := range privateQueries {
		fmt.Printf("Query Name: %s, Query ID: %s\n", *query.Name, query.Id.String())
	}
}

func (p *Ado) getProjectID(projectName string) (*core.TeamProject, error) {
	ctx := context.Background()
	coreClient, err := core.NewClient(p.ctx, p.connection)
	if err != nil {
		return nil, err
	}

	project, err := coreClient.GetProject(ctx, core.GetProjectArgs{
		ProjectId: &projectName,
	})
	if err != nil {
		return nil, err
	}

	return project, nil
}

func getQueries(items *[]workitemtracking.QueryHierarchyItem) []workitemtracking.QueryHierarchyItem {

	privateQueries := []workitemtracking.QueryHierarchyItem{}
	newQueryLevel := []uuid.UUID{}

	for _, query := range *items {
		if *query.IsPublic {
			continue
		}

		if query.IsFolder != nil && *query.IsFolder {
			if !(query.HasChildren != nil && *query.HasChildren) {
				continue
			}

			if query.Children == nil {
				newQueryLevel = append(newQueryLevel, *query.Id)
			} else {

				c := getQueries(query.Children)
				privateQueries = slices.Concat(privateQueries, c)
			}
		} else {
			privateQueries = append(privateQueries, query)
			// fmt.Println(*query.Path, *query.Name, *query.Wiql, '\n')
		}

		// fmt.Println(*query.HasChildren)
	}

	l := len(newQueryLevel)
	fmt.Println(">>>", strconv.FormatInt(int64(l), 10))
	// newQueryLevel
	// b := &[]uuid.UUID{*query.Id}
	// myq, err := workItemClient.GetQueriesBatch(p.ctx, workitemtracking.GetQueriesBatchArgs{
	// 	QueryGetRequest: &workitemtracking.QueryBatchGetRequest{
	// 		Expand: &a,
	// 		// ErrorPolicy: &"Fail",
	// 		Ids: &newQueryLevel,
	// 	},
	// 	Project: &projectID,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(myq)

	return privateQueries
}

func (p *Ado) listPrivateQueries(projectID string) ([]workitemtracking.QueryHierarchyItem, error) {
	workItemClient, err := workitemtracking.NewClient(p.ctx, p.connection)
	if err != nil {
		return nil, err
	}

	a := workitemtracking.QueryExpandValues.All
	depth := 2
	queries, err := workItemClient.GetQueries(p.ctx, workitemtracking.GetQueriesArgs{
		Project: &projectID,
		Expand:  &a,
		Depth:   &depth,
	})
	if err != nil {
		return nil, err
	}

	pq := getQueries(queries)

	// privateQueries := []workitemtracking.QueryHierarchyItem{}
	// for _, query := range *queries {
	// 	if !*query.IsPublic {

	// 		fmt.Println(*query.HasChildren)
	// 		privateQueries = append(privateQueries, query)

	// 	}
	// }

	privateQueries := pq

	return privateQueries, nil
}

func (p *Ado) CreateIteration() {
	project := "DevOps"
	team := "Core-Fusion Team"

	// q,_:=tokenadmin.NewClient(p.ctx,p.connection)
	// q.ListPersonalAccessTokens(p.ctx,tokenadmin.ListPersonalAccessTokensArgs{})

	// coreClient, _ := core.NewClient(p.ctx, p.connection)
	// proj, _ := coreClient.GetProject(p.ctx, core.GetProjectArgs{
	// 	ProjectId: &project,
	// })

	workClient, _ := work.NewClient(p.ctx, p.connection)

	// startDate := azuredevops.Time{Time: time.Now()}
	// finishDate := azuredevops.Time{Time: startDate.Time.AddDate(0, 0, 14)} // 2 weeks iteration
	// iterationName := work.TimeFrame{TimeFrame: "New Iteration"}

	// iteration := work.TeamIterationAttributes{
	// 	FinishDate: &finishDate,
	// 	StartDate:  &startDate,
	// 	TimeFrame:  &iterationName,
	// }

	// workClient.PostTeamIteration(p.ctx, work.PostTeamIterationArgs{
	// 	Iteration: &iteration,
	// 	Project:   new(string),
	// 	Team:      new(string),
	// })

	//Get list of iteration added to a team
	a, _ := workClient.GetTeamIterations(p.ctx, work.GetTeamIterationsArgs{
		Project: &project,
		Team:    &team,
	})

	fmt.Println(a)

	// workClient.PostTeamIteration(p.ctx, work.PostTeamIterationArgs{
	// 	Iteration: &work.TeamSettingsIteration{},
	// 	Project:   new(string),
	// 	Team:      new(string),
	// })

}

func (p *Ado) Test() {

	c, _ := graph.NewClient(p.ctx, p.connection)
	c.GetUser(p.ctx, graph.GetUserArgs{
		UserDescriptor: new(string),
	})

	coreClient, _ := core.NewClient(p.ctx, p.connection)
	projects, _ := coreClient.GetProjects(p.ctx, core.GetProjectsArgs{})

	for _, project := range projects.Value {
		fmt.Println(*project.Name, project.Id.String())

		pid := project.Id.String()
		teams, _ := coreClient.GetTeams(p.ctx, core.GetTeamsArgs{
			ProjectId: &pid,
			// ExpandIdentity: new(bool),
		})

		for _, team := range *teams {
			fmt.Println("  > ", *team.Name)

			tid := team.Id.String()
			members, _ := coreClient.GetTeamMembersWithExtendedProperties(p.ctx, core.GetTeamMembersWithExtendedPropertiesArgs{
				ProjectId: &pid,
				TeamId:    &tid,
			})

			for _, member := range *members {
				if member.IsTeamAdmin != nil {
					fmt.Println("    > ", *member.Identity.DisplayName, *member.Identity.Url, member.Identity.Links, member.Identity.Id, "Admin")
				} else {
					fmt.Println("    > ", *member.Identity.DisplayName, *member.Identity.Url, member.Identity.Links, member.Identity.Id)
				}

			}
		}

	}

}

// coreClient.GetTeams(ctx, core.GetTeamsArgs{})

// fmt.Println("==> ", projects.Value)

// s := security.NewClient(ctx, connection)
// s.
// c:=cix.NewClient(ctx,connection)
// a,_:=c.GetConfigurations(ctx,cix.GetConfigurationsArgs{})

// c,_:=accounts.NewClient(ctx,connection)
// a,_ :=c.GetAccounts(ctx,accounts.GetAccountsArgs{})

// c:=operations.NewClient(ctx,connection)
// c.GetOperation(ctx,operations.GetOperationArgs{})

// setClient:= settings.NewClient(ctx, connection)
// setClient.GetEntries(ctx, settings.GetEntriesArgs{})

// workClient, _ := work.NewClient(ctx, connection)

// fmt.Println(connection.AuthorizationString)
