/*
 * Formance Stack API
 * Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 
 *
 * The version of the OpenAPI document: develop
 * Contact: support@formance.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.formance.formance.model;

import java.util.Objects;
import java.util.Arrays;
import com.formance.formance.model.Posting;
import com.formance.formance.model.Volume;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import java.io.IOException;
import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Transaction
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class Transaction {
  public static final String SERIALIZED_NAME_TIMESTAMP = "timestamp";
  @SerializedName(SERIALIZED_NAME_TIMESTAMP)
  private OffsetDateTime timestamp;

  public static final String SERIALIZED_NAME_POSTINGS = "postings";
  @SerializedName(SERIALIZED_NAME_POSTINGS)
  private List<Posting> postings = new ArrayList<>();

  public static final String SERIALIZED_NAME_REFERENCE = "reference";
  @SerializedName(SERIALIZED_NAME_REFERENCE)
  private String reference;

  public static final String SERIALIZED_NAME_METADATA = "metadata";
  @SerializedName(SERIALIZED_NAME_METADATA)
  private Map<String, Object> metadata = null;

  public static final String SERIALIZED_NAME_TXID = "txid";
  @SerializedName(SERIALIZED_NAME_TXID)
  private Long txid;

  public static final String SERIALIZED_NAME_PRE_COMMIT_VOLUMES = "preCommitVolumes";
  @SerializedName(SERIALIZED_NAME_PRE_COMMIT_VOLUMES)
  private Map<String, Map<String, Volume>> preCommitVolumes = null;

  public static final String SERIALIZED_NAME_POST_COMMIT_VOLUMES = "postCommitVolumes";
  @SerializedName(SERIALIZED_NAME_POST_COMMIT_VOLUMES)
  private Map<String, Map<String, Volume>> postCommitVolumes = null;

  public Transaction() {
  }

  public Transaction timestamp(OffsetDateTime timestamp) {
    
    this.timestamp = timestamp;
    return this;
  }

   /**
   * Get timestamp
   * @return timestamp
  **/
  @javax.annotation.Nonnull

  public OffsetDateTime getTimestamp() {
    return timestamp;
  }


  public void setTimestamp(OffsetDateTime timestamp) {
    this.timestamp = timestamp;
  }


  public Transaction postings(List<Posting> postings) {
    
    this.postings = postings;
    return this;
  }

  public Transaction addPostingsItem(Posting postingsItem) {
    this.postings.add(postingsItem);
    return this;
  }

   /**
   * Get postings
   * @return postings
  **/
  @javax.annotation.Nonnull

  public List<Posting> getPostings() {
    return postings;
  }


  public void setPostings(List<Posting> postings) {
    this.postings = postings;
  }


  public Transaction reference(String reference) {
    
    this.reference = reference;
    return this;
  }

   /**
   * Get reference
   * @return reference
  **/
  @javax.annotation.Nullable

  public String getReference() {
    return reference;
  }


  public void setReference(String reference) {
    this.reference = reference;
  }


  public Transaction metadata(Map<String, Object> metadata) {
    
    this.metadata = metadata;
    return this;
  }

  public Transaction putMetadataItem(String key, Object metadataItem) {
    if (this.metadata == null) {
      this.metadata = new HashMap<>();
    }
    this.metadata.put(key, metadataItem);
    return this;
  }

   /**
   * Get metadata
   * @return metadata
  **/
  @javax.annotation.Nullable

  public Map<String, Object> getMetadata() {
    return metadata;
  }


  public void setMetadata(Map<String, Object> metadata) {
    this.metadata = metadata;
  }


  public Transaction txid(Long txid) {
    
    this.txid = txid;
    return this;
  }

   /**
   * Get txid
   * minimum: 0
   * @return txid
  **/
  @javax.annotation.Nonnull

  public Long getTxid() {
    return txid;
  }


  public void setTxid(Long txid) {
    this.txid = txid;
  }


  public Transaction preCommitVolumes(Map<String, Map<String, Volume>> preCommitVolumes) {
    
    this.preCommitVolumes = preCommitVolumes;
    return this;
  }

  public Transaction putPreCommitVolumesItem(String key, Map<String, Volume> preCommitVolumesItem) {
    if (this.preCommitVolumes == null) {
      this.preCommitVolumes = new HashMap<>();
    }
    this.preCommitVolumes.put(key, preCommitVolumesItem);
    return this;
  }

   /**
   * Get preCommitVolumes
   * @return preCommitVolumes
  **/
  @javax.annotation.Nullable

  public Map<String, Map<String, Volume>> getPreCommitVolumes() {
    return preCommitVolumes;
  }


  public void setPreCommitVolumes(Map<String, Map<String, Volume>> preCommitVolumes) {
    this.preCommitVolumes = preCommitVolumes;
  }


  public Transaction postCommitVolumes(Map<String, Map<String, Volume>> postCommitVolumes) {
    
    this.postCommitVolumes = postCommitVolumes;
    return this;
  }

  public Transaction putPostCommitVolumesItem(String key, Map<String, Volume> postCommitVolumesItem) {
    if (this.postCommitVolumes == null) {
      this.postCommitVolumes = new HashMap<>();
    }
    this.postCommitVolumes.put(key, postCommitVolumesItem);
    return this;
  }

   /**
   * Get postCommitVolumes
   * @return postCommitVolumes
  **/
  @javax.annotation.Nullable

  public Map<String, Map<String, Volume>> getPostCommitVolumes() {
    return postCommitVolumes;
  }


  public void setPostCommitVolumes(Map<String, Map<String, Volume>> postCommitVolumes) {
    this.postCommitVolumes = postCommitVolumes;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    Transaction transaction = (Transaction) o;
    return Objects.equals(this.timestamp, transaction.timestamp) &&
        Objects.equals(this.postings, transaction.postings) &&
        Objects.equals(this.reference, transaction.reference) &&
        Objects.equals(this.metadata, transaction.metadata) &&
        Objects.equals(this.txid, transaction.txid) &&
        Objects.equals(this.preCommitVolumes, transaction.preCommitVolumes) &&
        Objects.equals(this.postCommitVolumes, transaction.postCommitVolumes);
  }

  @Override
  public int hashCode() {
    return Objects.hash(timestamp, postings, reference, metadata, txid, preCommitVolumes, postCommitVolumes);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class Transaction {\n");
    sb.append("    timestamp: ").append(toIndentedString(timestamp)).append("\n");
    sb.append("    postings: ").append(toIndentedString(postings)).append("\n");
    sb.append("    reference: ").append(toIndentedString(reference)).append("\n");
    sb.append("    metadata: ").append(toIndentedString(metadata)).append("\n");
    sb.append("    txid: ").append(toIndentedString(txid)).append("\n");
    sb.append("    preCommitVolumes: ").append(toIndentedString(preCommitVolumes)).append("\n");
    sb.append("    postCommitVolumes: ").append(toIndentedString(postCommitVolumes)).append("\n");
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
