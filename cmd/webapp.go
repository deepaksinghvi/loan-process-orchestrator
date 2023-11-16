/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	controller "github.com/deepaksinghvi/loan-process-orchestrator/api"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/spf13/cobra"
)

// webappCmd represents the webapp command
var webappCmd = &cobra.Command{
	Use:   "webapp",
	Short: "Loan Origination Web App",
	Long:  `Loan Origination Worker App to interact with Temporal Workflow Engine Worker.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Loan Origination webapp")
		starter()
	},
}

func starter() {
	router := setupRouter()
	fmt.Println("Loan Origination loan_worker started, press ctrl+c to terminate...")
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	// Swagger setup
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Loan app apis
	router.POST("/loan-application", controller.CreateLoanApplication)
	router.GET("/loan-application/:workflow_id/:run_id", controller.GetLoanApplication)
	router.POST("/loan-application-approval/:workflow_id/:run_id", controller.LoanApplicationApproval)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	return router
}
func init() {
	rootCmd.AddCommand(webappCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webappCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webappCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
