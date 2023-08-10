<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;


class InsertConfigResponse
{
    /**
     * Config created successfully.
     * 
     * @var ?\formance\stack\Models\Shared\ConfigResponse $configResponse
     */
	
    public ?\formance\stack\Models\Shared\ConfigResponse $configResponse = null;
    
	
    public string $contentType;
    
	
    public int $statusCode;
    
	
    public ?\Psr\Http\Message\ResponseInterface $rawResponse = null;
    
    /**
     * Error
     * 
     * @var ?\formance\stack\Models\Shared\WebhooksErrorResponse $webhooksErrorResponse
     */
	
    public ?\formance\stack\Models\Shared\WebhooksErrorResponse $webhooksErrorResponse = null;
    
	public function __construct()
	{
		$this->configResponse = null;
		$this->contentType = "";
		$this->statusCode = 0;
		$this->rawResponse = null;
		$this->webhooksErrorResponse = null;
	}
}
