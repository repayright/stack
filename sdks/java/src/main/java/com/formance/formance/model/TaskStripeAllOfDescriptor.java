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
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import java.io.IOException;

/**
 * TaskStripeAllOfDescriptor
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class TaskStripeAllOfDescriptor {
  public static final String SERIALIZED_NAME_NAME = "name";
  @SerializedName(SERIALIZED_NAME_NAME)
  private String name;

  public static final String SERIALIZED_NAME_MAIN = "main";
  @SerializedName(SERIALIZED_NAME_MAIN)
  private Boolean main;

  public static final String SERIALIZED_NAME_ACCOUNT = "account";
  @SerializedName(SERIALIZED_NAME_ACCOUNT)
  private String account;

  public TaskStripeAllOfDescriptor() {
  }

  public TaskStripeAllOfDescriptor name(String name) {
    
    this.name = name;
    return this;
  }

   /**
   * Get name
   * @return name
  **/
  @javax.annotation.Nonnull

  public String getName() {
    return name;
  }


  public void setName(String name) {
    this.name = name;
  }


  public TaskStripeAllOfDescriptor main(Boolean main) {
    
    this.main = main;
    return this;
  }

   /**
   * Get main
   * @return main
  **/
  @javax.annotation.Nullable

  public Boolean getMain() {
    return main;
  }


  public void setMain(Boolean main) {
    this.main = main;
  }


  public TaskStripeAllOfDescriptor account(String account) {
    
    this.account = account;
    return this;
  }

   /**
   * Get account
   * @return account
  **/
  @javax.annotation.Nonnull

  public String getAccount() {
    return account;
  }


  public void setAccount(String account) {
    this.account = account;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    TaskStripeAllOfDescriptor taskStripeAllOfDescriptor = (TaskStripeAllOfDescriptor) o;
    return Objects.equals(this.name, taskStripeAllOfDescriptor.name) &&
        Objects.equals(this.main, taskStripeAllOfDescriptor.main) &&
        Objects.equals(this.account, taskStripeAllOfDescriptor.account);
  }

  @Override
  public int hashCode() {
    return Objects.hash(name, main, account);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class TaskStripeAllOfDescriptor {\n");
    sb.append("    name: ").append(toIndentedString(name)).append("\n");
    sb.append("    main: ").append(toIndentedString(main)).append("\n");
    sb.append("    account: ").append(toIndentedString(account)).append("\n");
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

