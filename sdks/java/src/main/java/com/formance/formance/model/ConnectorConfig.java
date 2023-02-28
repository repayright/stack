/*
 * Formance Stack API
 * Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 
 *
 * The version of the OpenAPI document: v1.0.20230228
 * Contact: support@formance.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.formance.formance.model;

import java.util.Objects;
import java.util.Arrays;
import com.formance.formance.model.BankingCircleConfig;
import com.formance.formance.model.CurrencyCloudConfig;
import com.formance.formance.model.DummyPayConfig;
import com.formance.formance.model.ModulrConfig;
import com.formance.formance.model.StripeConfig;
import com.formance.formance.model.WiseConfig;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import java.io.IOException;

/**
 * ConnectorConfig
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class ConnectorConfig {
  public static final String SERIALIZED_NAME_POLLING_PERIOD = "pollingPeriod";
  @SerializedName(SERIALIZED_NAME_POLLING_PERIOD)
  private String pollingPeriod;

  public static final String SERIALIZED_NAME_API_KEY = "apiKey";
  @SerializedName(SERIALIZED_NAME_API_KEY)
  private String apiKey;

  public static final String SERIALIZED_NAME_PAGE_SIZE = "pageSize";
  @SerializedName(SERIALIZED_NAME_PAGE_SIZE)
  private Long pageSize = 10l;

  public static final String SERIALIZED_NAME_FILE_POLLING_PERIOD = "filePollingPeriod";
  @SerializedName(SERIALIZED_NAME_FILE_POLLING_PERIOD)
  private String filePollingPeriod = "10s";

  public static final String SERIALIZED_NAME_FILE_GENERATION_PERIOD = "fileGenerationPeriod";
  @SerializedName(SERIALIZED_NAME_FILE_GENERATION_PERIOD)
  private String fileGenerationPeriod = "10s";

  public static final String SERIALIZED_NAME_DIRECTORY = "directory";
  @SerializedName(SERIALIZED_NAME_DIRECTORY)
  private String directory;

  public static final String SERIALIZED_NAME_API_SECRET = "apiSecret";
  @SerializedName(SERIALIZED_NAME_API_SECRET)
  private String apiSecret;

  public static final String SERIALIZED_NAME_ENDPOINT = "endpoint";
  @SerializedName(SERIALIZED_NAME_ENDPOINT)
  private String endpoint;

  public static final String SERIALIZED_NAME_LOGIN_I_D = "loginID";
  @SerializedName(SERIALIZED_NAME_LOGIN_I_D)
  private String loginID;

  public static final String SERIALIZED_NAME_USERNAME = "username";
  @SerializedName(SERIALIZED_NAME_USERNAME)
  private String username;

  public static final String SERIALIZED_NAME_PASSWORD = "password";
  @SerializedName(SERIALIZED_NAME_PASSWORD)
  private String password;

  public static final String SERIALIZED_NAME_AUTHORIZATION_ENDPOINT = "authorizationEndpoint";
  @SerializedName(SERIALIZED_NAME_AUTHORIZATION_ENDPOINT)
  private String authorizationEndpoint;

  public ConnectorConfig() {
  }

  public ConnectorConfig pollingPeriod(String pollingPeriod) {
    
    this.pollingPeriod = pollingPeriod;
    return this;
  }

   /**
   * The frequency at which the connector will fetch transactions
   * @return pollingPeriod
  **/
  @javax.annotation.Nullable

  public String getPollingPeriod() {
    return pollingPeriod;
  }


  public void setPollingPeriod(String pollingPeriod) {
    this.pollingPeriod = pollingPeriod;
  }


  public ConnectorConfig apiKey(String apiKey) {
    
    this.apiKey = apiKey;
    return this;
  }

   /**
   * Get apiKey
   * @return apiKey
  **/
  @javax.annotation.Nonnull

  public String getApiKey() {
    return apiKey;
  }


  public void setApiKey(String apiKey) {
    this.apiKey = apiKey;
  }


  public ConnectorConfig pageSize(Long pageSize) {
    
    this.pageSize = pageSize;
    return this;
  }

   /**
   * Number of BalanceTransaction to fetch at each polling interval. 
   * minimum: 0
   * @return pageSize
  **/
  @javax.annotation.Nullable

  public Long getPageSize() {
    return pageSize;
  }


  public void setPageSize(Long pageSize) {
    this.pageSize = pageSize;
  }


  public ConnectorConfig filePollingPeriod(String filePollingPeriod) {
    
    this.filePollingPeriod = filePollingPeriod;
    return this;
  }

   /**
   * The frequency at which the connector will try to fetch new payment objects from the directory
   * @return filePollingPeriod
  **/
  @javax.annotation.Nullable

  public String getFilePollingPeriod() {
    return filePollingPeriod;
  }


  public void setFilePollingPeriod(String filePollingPeriod) {
    this.filePollingPeriod = filePollingPeriod;
  }


  public ConnectorConfig fileGenerationPeriod(String fileGenerationPeriod) {
    
    this.fileGenerationPeriod = fileGenerationPeriod;
    return this;
  }

   /**
   * The frequency at which the connector will create new payment objects in the directory
   * @return fileGenerationPeriod
  **/
  @javax.annotation.Nullable

  public String getFileGenerationPeriod() {
    return fileGenerationPeriod;
  }


  public void setFileGenerationPeriod(String fileGenerationPeriod) {
    this.fileGenerationPeriod = fileGenerationPeriod;
  }


  public ConnectorConfig directory(String directory) {
    
    this.directory = directory;
    return this;
  }

   /**
   * Get directory
   * @return directory
  **/
  @javax.annotation.Nonnull

  public String getDirectory() {
    return directory;
  }


  public void setDirectory(String directory) {
    this.directory = directory;
  }


  public ConnectorConfig apiSecret(String apiSecret) {
    
    this.apiSecret = apiSecret;
    return this;
  }

   /**
   * Get apiSecret
   * @return apiSecret
  **/
  @javax.annotation.Nonnull

  public String getApiSecret() {
    return apiSecret;
  }


  public void setApiSecret(String apiSecret) {
    this.apiSecret = apiSecret;
  }


  public ConnectorConfig endpoint(String endpoint) {
    
    this.endpoint = endpoint;
    return this;
  }

   /**
   * Get endpoint
   * @return endpoint
  **/
  @javax.annotation.Nonnull

  public String getEndpoint() {
    return endpoint;
  }


  public void setEndpoint(String endpoint) {
    this.endpoint = endpoint;
  }


  public ConnectorConfig loginID(String loginID) {
    
    this.loginID = loginID;
    return this;
  }

   /**
   * Username of the API Key holder
   * @return loginID
  **/
  @javax.annotation.Nonnull

  public String getLoginID() {
    return loginID;
  }


  public void setLoginID(String loginID) {
    this.loginID = loginID;
  }


  public ConnectorConfig username(String username) {
    
    this.username = username;
    return this;
  }

   /**
   * Get username
   * @return username
  **/
  @javax.annotation.Nonnull

  public String getUsername() {
    return username;
  }


  public void setUsername(String username) {
    this.username = username;
  }


  public ConnectorConfig password(String password) {
    
    this.password = password;
    return this;
  }

   /**
   * Get password
   * @return password
  **/
  @javax.annotation.Nonnull

  public String getPassword() {
    return password;
  }


  public void setPassword(String password) {
    this.password = password;
  }


  public ConnectorConfig authorizationEndpoint(String authorizationEndpoint) {
    
    this.authorizationEndpoint = authorizationEndpoint;
    return this;
  }

   /**
   * Get authorizationEndpoint
   * @return authorizationEndpoint
  **/
  @javax.annotation.Nonnull

  public String getAuthorizationEndpoint() {
    return authorizationEndpoint;
  }


  public void setAuthorizationEndpoint(String authorizationEndpoint) {
    this.authorizationEndpoint = authorizationEndpoint;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    ConnectorConfig connectorConfig = (ConnectorConfig) o;
    return Objects.equals(this.pollingPeriod, connectorConfig.pollingPeriod) &&
        Objects.equals(this.apiKey, connectorConfig.apiKey) &&
        Objects.equals(this.pageSize, connectorConfig.pageSize) &&
        Objects.equals(this.filePollingPeriod, connectorConfig.filePollingPeriod) &&
        Objects.equals(this.fileGenerationPeriod, connectorConfig.fileGenerationPeriod) &&
        Objects.equals(this.directory, connectorConfig.directory) &&
        Objects.equals(this.apiSecret, connectorConfig.apiSecret) &&
        Objects.equals(this.endpoint, connectorConfig.endpoint) &&
        Objects.equals(this.loginID, connectorConfig.loginID) &&
        Objects.equals(this.username, connectorConfig.username) &&
        Objects.equals(this.password, connectorConfig.password) &&
        Objects.equals(this.authorizationEndpoint, connectorConfig.authorizationEndpoint);
  }

  @Override
  public int hashCode() {
    return Objects.hash(pollingPeriod, apiKey, pageSize, filePollingPeriod, fileGenerationPeriod, directory, apiSecret, endpoint, loginID, username, password, authorizationEndpoint);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ConnectorConfig {\n");
    sb.append("    pollingPeriod: ").append(toIndentedString(pollingPeriod)).append("\n");
    sb.append("    apiKey: ").append(toIndentedString(apiKey)).append("\n");
    sb.append("    pageSize: ").append(toIndentedString(pageSize)).append("\n");
    sb.append("    filePollingPeriod: ").append(toIndentedString(filePollingPeriod)).append("\n");
    sb.append("    fileGenerationPeriod: ").append(toIndentedString(fileGenerationPeriod)).append("\n");
    sb.append("    directory: ").append(toIndentedString(directory)).append("\n");
    sb.append("    apiSecret: ").append(toIndentedString(apiSecret)).append("\n");
    sb.append("    endpoint: ").append(toIndentedString(endpoint)).append("\n");
    sb.append("    loginID: ").append(toIndentedString(loginID)).append("\n");
    sb.append("    username: ").append(toIndentedString(username)).append("\n");
    sb.append("    password: ").append(toIndentedString(password)).append("\n");
    sb.append("    authorizationEndpoint: ").append(toIndentedString(authorizationEndpoint)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

