/////////////////////////////////////////////////////////////////////////////
//
//  @file efml.h
//
//  @brief API interface of Enflame Management Library
//
//  Enflame Tech, All Rights Reserved. 2019-2025 Copyright (C)
//
/////////////////////////////////////////////////////////////////////////////

#if !defined(_EFML_H_)
#define _EFML_H_

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
* @brief Enflame data structure for device thermal info.
*
*/
typedef struct {
   float cur_dev_temp;
   float cur_hbm0_temp;
   float cur_hbm1_temp;
}efmlDevThermalInfo_t;

 /**
* @brief Enflame data structure for DTU usage info.
*
*/
typedef struct {
        efmlHbmUsage_t hbm;
        efmlSipUsage_t sip;
        efmlCqmUsage_t cqm;
}efmlDtuUsage_t;

typedef enum {
    EFML_ESL_LINK_SPEED_GEN1      = 1,
    EFML_ESL_LINK_SPEED_GEN2      = 2,
    EFML_ESL_LINK_SPEED_GEN3      = 3,
    EFML_ESL_LINK_SPEED_GEN4      = 4,
    EFML_ESL_LINK_SPEED_GEN5      = 5,
    EFML_ESL_LINK_SPEED_ESM_2P5GT = 0x10,
    EFML_ESL_LINK_SPEED_ESM_5GT,
    EFML_ESL_LINK_SPEED_ESM_8GT,
    EFML_ESL_LINK_SPEED_ESM_16GT,
    EFML_ESL_LINK_SPEED_ESM_20GT,
    EFML_ESL_LINK_SPEED_ESM_25GT
} efmlEslSpeed_t;

typedef enum {
    EFML_ESL_LINK_WIDTH_X1  = 1,
    EFML_ESL_LINK_WIDTH_X2  = 2,
    EFML_ESL_LINK_WIDTH_X4  = 4,
    EFML_ESL_LINK_WIDTH_X8  = 8,
    EFML_ESL_LINK_WIDTH_X16 = 16,
} efmlEslWidth_t;

typedef enum {
    EFML_ESL_PORT_RC = 0,
    EFML_ESL_PORT_EP = 1,
} efmlEslPortType_t;

/**
* @brief Enflame data structure for esl link info.
*
*/
typedef struct {
    efmlEslSpeed_t link_speed;
    efmlEslSpeed_t max_link_speed;
    efmlEslWidth_t link_width;
    efmlEslWidth_t max_link_width;
} efmlEslLinkInfo_t;

/**
*@brief Enflame data structure for Esl Throughput info.
*
*/
typedef struct {
    float    tx_throughput;
    float    rx_throughput;
    uint64_t tx_nak;
    uint64_t rx_nak;
} efmlEslThroughputInfo_t;


/**
 * @brief Enflame data structure for esl port info.
 *
 */
typedef struct {
    uint32_t          connected;
    char              uuid[16];
    uint32_t          vendor_id;
    uint32_t          device_id;
    uint32_t          domain_id;
    uint32_t          bus_id;
    uint32_t          dev_id;
    uint32_t          func_id;
    uint32_t          port_id;
    efmlEslPortType_t port_type;

    uint32_t          remote_card_id;
    char              remote_uuid[16];
    uint32_t          remote_vendor_id;
    uint32_t          remote_device_id;
    uint32_t          remote_domain_id;
    uint32_t          remote_bus_id;
    uint32_t          remote_dev_id;
    uint32_t          remote_func_id;
    uint32_t          remote_port_id;
    efmlEslPortType_t remote_port_type;
} efmlEslPortInfo_t;

/**
 * @brief Enflame Management Library initialization with log file name
 *
 * @return efmlReturn_t
 */
efmlReturn_t EfmlInitFile(char *log_file_name, bool no_driver);

/**
 * @brief Enflame Management Library get the device temperature.
 *
 */
efmlReturn_t EfmlGetDevTemp(uint32_t dev_idx, efmlDevThermalInfo_t* p_temp);

/**
 * @brief Enflame Management Library get the device DPM level.
 *
 */
efmlReturn_t EfmlGetDevDpmLevel(uint32_t dev_idx, uint32_t* p_dpm);

/**
 * @brief Enflame Management Library switch dtu power stock mode.
 *
 */
efmlReturn_t EfmlSetDevSupportPowerStock(uint32_t dev_idx, bool enable_power_stock_support);

/**
 * @brief Enflame Management Library get if dtu support power stock mode.
 *
 */
efmlReturn_t EfmlGetDevSupportPowerStock(uint32_t dev_idx, bool* is_support_power_stock);

/**
 * @brief Enflame Management Library dump supported device list.
 *
 */
efmlReturn_t EfmlDumpDevList(void);

/**
 * @brief Enflame Management Library get device parent info.
 *
 */
efmlReturn_t EfmlGetDevParentInfo(uint32_t dev_idx, efmlDeviceInfo_t *p_info);


/**
 * @brief Enflame Management Library get device PCIe slot number.
 *
 */
efmlReturn_t EfmlGetDevSlotNum(uint32_t dev_idx, uint32_t *p_slot);

/**
 * @brief Enflame Management Library select one target device by index.
 *
 */
efmlReturn_t EfmlSelDevByIndex(uint32_t dev_idx);

/**
 * @brief Enflame Management Library get pcie up top.
 *
 */
efmlReturn_t EfmlGetDevTop(uint32_t dev_idx);

/**
 * @brief Enflame Management Library get virtualization status.
 *
 */
efmlReturn_t EfmlGetDevIsVdtuEnabled(uint32_t dev_idx, bool *is_vdtu_enabled);

/**
 * @brief Enflame Management Library get if dtu support low power mode.
 *
 */
efmlReturn_t EfmlGetDevSupportLowPower(uint32_t dev_idx, bool *is_support_low_power);

/**
 * @brief Enflame Management Library switch dtu low power mode.
 *
 */
efmlReturn_t EfmlSetDevSupportLowPower(uint32_t dev_idx, bool enable_low_power_support);

/**
 * @brief Enflame Management Library get total esl port numbers.
 *
 */
efmlReturn_t EfmlGetEslPortNum(uint32_t dev_idx, uint32_t *p_data);

/**
 * @brief Enflame Management Library get esl port info.
 *
 */
efmlReturn_t EfmlGetEslPortInfo(uint32_t dev_idx, uint32_t port_id, efmlEslPortInfo_t *p_info);

/**
 * @brief Enflame Management Library get esl link info.
 *
 */
efmlReturn_t EfmlGetEslLinkInfo(uint32_t dev_idx, uint32_t port_id, efmlEslLinkInfo_t *p_info);

/**
 * @brief Enflame Management Library get esl dtuid info.
 *
 */
efmlReturn_t EfmlGetEslDtuId(uint32_t dev_idx, uint32_t *p_data);

/**
 * @brief Enflame Management Library get esl support info.
 *
 */
efmlReturn_t EfmlGetEslIsSupported(uint32_t dev_idx, bool *is_esl_supported);

/**
 * @brief Enflame Management Library get esl throughput info.
 *
 */
efmlReturn_t EfmlGetEslThroughput(
    uint32_t dev_idx, uint32_t port_id, efmlEslThroughputInfo_t *p_info);

#ifdef __cplusplus
} // extern "C"
#endif


#endif  //__EFML_H_
